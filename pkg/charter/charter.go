package charter

import (
	"context"
	"github.com/jenkins-x-plugins/jx-charter/pkg/apis/chart/v1alpha1"
	"github.com/jenkins-x-plugins/jx-charter/pkg/client/clientset/versioned"
	"github.com/jenkins-x-plugins/jx-charter/pkg/helmdecoder"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UpsertChartFromSecret upserts the Chart CRD from the given Secret if its a helm secret
func UpsertChartFromSecret(ctx context.Context, chartClient versioned.Interface, r *v1.Secret) error {
	if r == nil {
		return nil
	}
	release, err := helmdecoder.ConvertSecretToHelmRelease(r)
	if err != nil {
		log.Logger().Warnf("failed to decode Secret %s/%s due to %v\n", r.Namespace, r.Namespace, err.Error())
		return nil
	}
	if release == nil {
		return nil
	}

	name := release.Name
	if name == "" {
		name = r.Name
	}
	namespace := release.Namespace
	if namespace == "" {
		namespace = r.Namespace
	}
	ch := &v1alpha1.Chart{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.APIVersion,
			Kind:       v1alpha1.KindChart,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: r.Annotations,
			Labels:      r.Labels,
		},
	}

	if release.Chart != nil && release.Chart.Metadata != nil {
		ch.Spec.Metadata = *release.Chart.Metadata
	}
	if release.Info != nil {
		ch.Status = v1alpha1.ToChartStatus(release.Info)
	}

	_, err = UpsertChart(ctx, chartClient, ch)
	return err
}

// UpsertChart upserts the Chart resource
func UpsertChart(ctx context.Context, chartClient versioned.Interface, ch *v1alpha1.Chart) (*v1alpha1.Chart, error) {
	ns := ch.Namespace
	name := ch.Name
	if name == "" {
		log.Logger().Warnf("missing chart name")
	}
	chartInterface := chartClient.ChartV1alpha1().Charts(ns)

	r, err := chartInterface.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return r, errors.Wrapf(err, "failed to get Chart resource %s/%s", ns, name)
		}
		r, err = chartInterface.Create(ctx, ch, metav1.CreateOptions{})
		if err != nil {
			return r, errors.Wrapf(err, "failed to create Chart resource %s/%s", ns, name)
		}
		return r, nil
	}
	r.Name = name
	r.Namespace = ns
	r.Spec = ch.Spec
	if ch.Status != nil {
		r.Status = ch.Status
	}
	r, err = chartInterface.Update(ctx, r, metav1.UpdateOptions{})
	if err != nil {
		return r, errors.Wrapf(err, "failed to update Chart resource %s/%s", ns, name)
	}
	return r, nil
}

// DeleteChartFromSecret deletes the Chart CRD from the given Secret if its a helm secret
func DeleteChartFromSecret(ctx context.Context, chartClient versioned.Interface, r *v1.Secret) error {
	if r == nil {
		return nil
	}
	release, err := helmdecoder.ConvertSecretToHelmRelease(r)
	if err != nil {
		log.Logger().Warnf("failed to decode Secret %s/%s due to %v\n", r.Namespace, r.Namespace, err.Error())
		return nil
	}
	if release == nil {
		return nil
	}

	ns := release.Namespace
	name := release.Name

	err = chartClient.ChartV1alpha1().Charts(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return errors.Wrapf(err, "failed to delete Chart %s/%s", ns, name)
	}
	return nil
}

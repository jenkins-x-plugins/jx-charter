/*
Copyright 2020 The Jenkins X Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/jenkins-x-plugins/jx-charter/pkg/apis/chart/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCharts implements ChartInterface
type FakeCharts struct {
	Fake *FakeChartV1alpha1
	ns   string
}

var chartsResource = v1alpha1.SchemeGroupVersion.WithResource("charts")

var chartsKind = v1alpha1.SchemeGroupVersion.WithKind("Chart")

// Get takes name of the chart, and returns the corresponding chart object, and an error if there is any.
func (c *FakeCharts) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Chart, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(chartsResource, c.ns, name), &v1alpha1.Chart{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Chart), err
}

// List takes label and field selectors, and returns the list of Charts that match those selectors.
func (c *FakeCharts) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ChartList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(chartsResource, chartsKind, c.ns, opts), &v1alpha1.ChartList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ChartList{ListMeta: obj.(*v1alpha1.ChartList).ListMeta}
	for _, item := range obj.(*v1alpha1.ChartList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested charts.
func (c *FakeCharts) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(chartsResource, c.ns, opts))

}

// Create takes the representation of a chart and creates it.  Returns the server's representation of the chart, and an error, if there is any.
func (c *FakeCharts) Create(ctx context.Context, chart *v1alpha1.Chart, opts v1.CreateOptions) (result *v1alpha1.Chart, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(chartsResource, c.ns, chart), &v1alpha1.Chart{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Chart), err
}

// Update takes the representation of a chart and updates it. Returns the server's representation of the chart, and an error, if there is any.
func (c *FakeCharts) Update(ctx context.Context, chart *v1alpha1.Chart, opts v1.UpdateOptions) (result *v1alpha1.Chart, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(chartsResource, c.ns, chart), &v1alpha1.Chart{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Chart), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCharts) UpdateStatus(ctx context.Context, chart *v1alpha1.Chart, opts v1.UpdateOptions) (*v1alpha1.Chart, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(chartsResource, "status", c.ns, chart), &v1alpha1.Chart{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Chart), err
}

// Delete takes name of the chart and deletes it. Returns an error if one occurs.
func (c *FakeCharts) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(chartsResource, c.ns, name, opts), &v1alpha1.Chart{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCharts) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(chartsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ChartList{})
	return err
}

// Patch applies the patch and returns the patched chart.
func (c *FakeCharts) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Chart, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(chartsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Chart{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Chart), err
}

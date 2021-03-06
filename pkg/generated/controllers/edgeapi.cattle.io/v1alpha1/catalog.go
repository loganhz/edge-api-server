/*
Copyright 2020 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type CatalogHandler func(string, *v1alpha1.Catalog) (*v1alpha1.Catalog, error)

type CatalogController interface {
	generic.ControllerMeta
	CatalogClient

	OnChange(ctx context.Context, name string, sync CatalogHandler)
	OnRemove(ctx context.Context, name string, sync CatalogHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() CatalogCache
}

type CatalogClient interface {
	Create(*v1alpha1.Catalog) (*v1alpha1.Catalog, error)
	Update(*v1alpha1.Catalog) (*v1alpha1.Catalog, error)
	UpdateStatus(*v1alpha1.Catalog) (*v1alpha1.Catalog, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Catalog, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.CatalogList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Catalog, err error)
}

type CatalogCache interface {
	Get(namespace, name string) (*v1alpha1.Catalog, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.Catalog, error)

	AddIndexer(indexName string, indexer CatalogIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Catalog, error)
}

type CatalogIndexer func(obj *v1alpha1.Catalog) ([]string, error)

type catalogController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewCatalogController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) CatalogController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &catalogController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromCatalogHandlerToHandler(sync CatalogHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Catalog
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Catalog))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *catalogController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Catalog))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateCatalogDeepCopyOnChange(client CatalogClient, obj *v1alpha1.Catalog, handler func(obj *v1alpha1.Catalog) (*v1alpha1.Catalog, error)) (*v1alpha1.Catalog, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *catalogController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *catalogController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *catalogController) OnChange(ctx context.Context, name string, sync CatalogHandler) {
	c.AddGenericHandler(ctx, name, FromCatalogHandlerToHandler(sync))
}

func (c *catalogController) OnRemove(ctx context.Context, name string, sync CatalogHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromCatalogHandlerToHandler(sync)))
}

func (c *catalogController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *catalogController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *catalogController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *catalogController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *catalogController) Cache() CatalogCache {
	return &catalogCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *catalogController) Create(obj *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	result := &v1alpha1.Catalog{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *catalogController) Update(obj *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	result := &v1alpha1.Catalog{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *catalogController) UpdateStatus(obj *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	result := &v1alpha1.Catalog{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *catalogController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *catalogController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Catalog, error) {
	result := &v1alpha1.Catalog{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *catalogController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.CatalogList, error) {
	result := &v1alpha1.CatalogList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *catalogController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *catalogController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1alpha1.Catalog, error) {
	result := &v1alpha1.Catalog{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type catalogCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *catalogCache) Get(namespace, name string) (*v1alpha1.Catalog, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1alpha1.Catalog), nil
}

func (c *catalogCache) List(namespace string, selector labels.Selector) (ret []*v1alpha1.Catalog, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Catalog))
	})

	return ret, err
}

func (c *catalogCache) AddIndexer(indexName string, indexer CatalogIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Catalog))
		},
	}))
}

func (c *catalogCache) GetByIndex(indexName, key string) (result []*v1alpha1.Catalog, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.Catalog, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Catalog))
	}
	return result, nil
}

type CatalogStatusHandler func(obj *v1alpha1.Catalog, status v1alpha1.CatalogStatus) (v1alpha1.CatalogStatus, error)

type CatalogGeneratingHandler func(obj *v1alpha1.Catalog, status v1alpha1.CatalogStatus) ([]runtime.Object, v1alpha1.CatalogStatus, error)

func RegisterCatalogStatusHandler(ctx context.Context, controller CatalogController, condition condition.Cond, name string, handler CatalogStatusHandler) {
	statusHandler := &catalogStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromCatalogHandlerToHandler(statusHandler.sync))
}

func RegisterCatalogGeneratingHandler(ctx context.Context, controller CatalogController, apply apply.Apply,
	condition condition.Cond, name string, handler CatalogGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &catalogGeneratingHandler{
		CatalogGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterCatalogStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type catalogStatusHandler struct {
	client    CatalogClient
	condition condition.Cond
	handler   CatalogStatusHandler
}

func (a *catalogStatusHandler) sync(key string, obj *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type catalogGeneratingHandler struct {
	CatalogGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *catalogGeneratingHandler) Remove(key string, obj *v1alpha1.Catalog) (*v1alpha1.Catalog, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.Catalog{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *catalogGeneratingHandler) Handle(obj *v1alpha1.Catalog, status v1alpha1.CatalogStatus) (v1alpha1.CatalogStatus, error) {
	objs, newStatus, err := a.CatalogGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}

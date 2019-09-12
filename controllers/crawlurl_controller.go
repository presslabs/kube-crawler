/*

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

package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	crawlv1 "github.com/presslabs/kube-crawler-controller/api/v1"
)

// CrawlURLReconciler reconciles a CrawlURL object
type CrawlURLReconciler struct {
	client.Client
	Log logr.Logger
}

const RecheckInterval = 30 * time.Second

// +kubebuilder:rbac:groups=crawl.presslabs.org,resources=crawlurls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crawl.presslabs.org,resources=crawlurls/status,verbs=get;update;patch

func (r *CrawlURLReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("crawlurl", req.NamespacedName)

	recheck := ctrl.Result{
		Requeue:      true,
		RequeueAfter: RecheckInterval,
	}

	obj := &crawlv1.CrawlURL{}
	if err := r.Get(ctx, req.NamespacedName, obj); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	// Skip crawling before RecheckInterval has passed
	lastCrawl := obj.Status.LastCrawlDate
	now := time.Now()
	if !now.After(lastCrawl.Add(RecheckInterval)) {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: lastCrawl.Add(RecheckInterval).Sub(now).Round(time.Second),
		}, nil
	}

	resp, err := http.Get(obj.Spec.URL)
	if err != nil {
		log.Error(err, "http error")
		return recheck, nil
	}
	log.Info("crawled url", "url", obj.Spec.URL, "status", resp.StatusCode)

	obj.Status.LastCrawlStatus = &resp.StatusCode
	obj.Status.LastCrawlDate = metav1.NewTime(now)

	if err := r.Status().Update(ctx, obj); err != nil {
		return ctrl.Result{}, err
	}

	return recheck, nil
}

func (r *CrawlURLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crawlv1.CrawlURL{}).
		Complete(r)
}

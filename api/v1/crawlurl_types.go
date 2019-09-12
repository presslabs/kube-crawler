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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CrawlURLSpec defines the desired state of CrawlURL
type CrawlURLSpec struct {
	URL string `json:"url"`
}

// CrawlURLStatus defines the observed state of CrawlURL
type CrawlURLStatus struct {
	LastCrawlDate   metav1.Time `json:"lastCrawlDate,omitempty"`
	LastCrawlStatus *int        `json:"lastCrawlStatus,omitempty"`
}

// +kubebuilder:object:root=true

// CrawlURL is the Schema for the crawlurls API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".spec.url",description="The URL to crawl"
// +kubebuilder:printcolumn:name="Status",type="integer",JSONPath=".status.lastCrawlStatus",description="Last crawl status"
// +kubebuilder:printcolumn:name="Last Crawl",type="date",JSONPath=".status.lastCrawlDate",description="Last crawl date"
type CrawlURL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CrawlURLSpec   `json:"spec,omitempty"`
	Status CrawlURLStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CrawlURLList contains a list of CrawlURL
type CrawlURLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CrawlURL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CrawlURL{}, &CrawlURLList{})
}

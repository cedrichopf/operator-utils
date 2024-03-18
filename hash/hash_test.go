package hash_test

import (
	"github.com/cedrichopf/operator-utils/hash"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test GenerateObjectHash", func() {
	Context("when object is valid", func() {
		It("returns the hash value", func() {
			object := map[string]string{
				"some": "value",
			}

			result, err := hash.GenerateObjectHash(&object)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("1916296322"))
		})
	})

	Context("when object is invalid", func() {
		It("returns the hash value", func() {
			object := make(chan string)

			result, err := hash.GenerateObjectHash(&object)
			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(""))
		})
	})
})

var _ = Describe("Test GenerateRevisionHashLabel", func() {
	Context("when function is called", func() {
		It("returns a map containing the value", func() {
			revisionHash := "test"

			result := hash.GenerateRevisionHashLabel(revisionHash)
			Expect(result[hash.REVISION_HASH_LABEL]).To(Equal("test"))
		})
	})
})

var _ = Describe("Test AddRevisionHashLabel", func() {
	Context("when object is valid", func() {
		It("adds the revision hash label to the object", func() {
			object := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "default",
				},
			}

			err := hash.AddRevisionHashLabel(object)
			Expect(err).ToNot(HaveOccurred())
			Expect(object.Labels[hash.REVISION_HASH_LABEL]).To(Equal("1356983499"))
		})
	})
})

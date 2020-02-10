/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
)

var _ = Describe("Testing CoherenceCluster", func() {

	It("cluster produces wka service name", func() {
		cluster := coherence.CoherenceCluster{}
		cluster.Name = "foo"

		Expect(cluster.GetWkaServiceName()).To(Equal("foo" + coherence.WKAServiceNameSuffix))
	})

	It("cluster has role with name", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		Expect(cluster.GetRole("storage")).To(Equal(roleOne))
	})

	It("cluster does not have role with name", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		Expect(cluster.GetRole("foo")).To(Equal(coherence.CoherenceRoleSpec{Role: "foo"}))
	})

	It("set cluster role", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		roleUpdate := coherence.CoherenceRoleSpec{}
		roleUpdate.Role = "storage"
		roleUpdate.Replicas = int32Ptr(19)

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		cluster.SetRole(roleUpdate)

		Expect(cluster.GetRole("storage")).To(Equal(roleUpdate))
	})

	It("set cluster role when role not in cluster", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		roleUpdate := coherence.CoherenceRoleSpec{}
		roleUpdate.Role = "foo"
		roleUpdate.Replicas = int32Ptr(19)

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		cluster.SetRole(roleUpdate)

		Expect(len(cluster.Spec.Roles)).To(Equal(2))
		Expect(cluster.GetRole("storage")).To(Equal(roleOne))
		Expect(cluster.GetRole("proxy")).To(Equal(roleTwo))
	})

	Context("loading CoherenceCluster from yaml file", func() {
		var (
			cluster coherence.CoherenceCluster
			file    []string
			err     error
		)

		JustBeforeEach(func() {
			cluster, err = helper.NewCoherenceClusterFromYaml("test-ns", file...)
		})

		When("file is valid", func() {
			BeforeEach(func() {
				file = []string{"test-coherence-cluster-one.yaml"}
			})

			It("should load fields", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cluster).NotTo(BeNil())

				// values come from the test-coherence-cluster-one.yaml file
				Expect(cluster.Name).To(Equal("test-cluster"))
				Expect(cluster.Spec.ReadinessProbe).NotTo(BeNil())
				Expect(cluster.Spec.ReadinessProbe.InitialDelaySeconds).To(Equal(int32Ptr(10)))
				Expect(cluster.Spec.ReadinessProbe.PeriodSeconds).To(Equal(int32Ptr(30)))

				Expect(cluster.Spec.Roles).ToNot(BeNil())
				Expect(len(cluster.Spec.Roles)).To(Equal(1))

				role := cluster.Spec.Roles[0]
				Expect(role.GetRoleName()).To(Equal("one"))
				Expect(role.GetReplicas()).To(Equal(int32(1)))
			})
		})

		When("file does not exist", func() {
			BeforeEach(func() {
				file = []string{"foo.yaml"}
			})

			It("should return error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		When("multiple yaml files", func() {
			BeforeEach(func() {
				file = []string{"test-coherence-cluster-one.yaml", "test-coherence-cluster-two.yaml"}
			})

			It("should load fields", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cluster).NotTo(BeNil())

				// values come from the test-coherence-cluster-one.yaml file and are then
				// overridden or added to by the test-coherence-cluster-two.yaml file
				Expect(cluster.Name).To(Equal("test-cluster-two"))
				Expect(cluster.Spec.ReadinessProbe).NotTo(BeNil())
				Expect(cluster.Spec.ReadinessProbe.InitialDelaySeconds).To(Equal(int32Ptr(60)))
				Expect(cluster.Spec.ReadinessProbe.PeriodSeconds).To(Equal(int32Ptr(30)))

				Expect(cluster.Spec.Roles).ToNot(BeNil())
				Expect(len(cluster.Spec.Roles)).To(Equal(2))

				roleOne := cluster.Spec.Roles[0]
				Expect(roleOne.GetRoleName()).To(Equal("one"))
				Expect(roleOne.GetReplicas()).To(Equal(int32(3)))

				roleTwo := cluster.Spec.Roles[1]
				Expect(roleTwo.GetRoleName()).To(Equal("two"))
				Expect(roleTwo.GetReplicas()).To(Equal(int32(3)))
			})
		})
	})
})

func TestDefaultCluster(t *testing.T) {
	g := NewGomegaWithT(t)
	cluster := coherence.CoherenceCluster{}
	g.Expect(cluster.GetClusterSize()).To(Equal(int(coherence.DefaultReplicas)))
}

func TestRoleWithDefaultReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	role := coherence.CoherenceRoleSpec{Role: "foo"}
	cluster := coherence.CoherenceCluster{
		Spec: coherence.CoherenceClusterSpec{Roles: []coherence.CoherenceRoleSpec{role}},
	}
	g.Expect(cluster.GetClusterSize()).To(Equal(int(coherence.DefaultReplicas)))
}

func TestRoleWithReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	role := coherence.CoherenceRoleSpec{Role: "foo"}
	role.SetReplicas(20)
	cluster := coherence.CoherenceCluster{
		Spec: coherence.CoherenceClusterSpec{Roles: []coherence.CoherenceRoleSpec{role}},
	}
	g.Expect(cluster.GetClusterSize()).To(Equal(20))
}

func TestTwoRoleWithDefaultReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	roleOne := coherence.CoherenceRoleSpec{Role: "foo"}
	roleTwo := coherence.CoherenceRoleSpec{Role: "bar"}
	cluster := coherence.CoherenceCluster{
		Spec: coherence.CoherenceClusterSpec{Roles: []coherence.CoherenceRoleSpec{roleOne, roleTwo}},
	}
	g.Expect(cluster.GetClusterSize()).To(Equal(6))
}

func TestTwoRoleWithSameReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	roleOne := coherence.CoherenceRoleSpec{Role: "foo"}
	roleOne.SetReplicas(1)
	roleTwo := coherence.CoherenceRoleSpec{Role: "bar"}
	roleTwo.SetReplicas(1)
	cluster := coherence.CoherenceCluster{
		Spec: coherence.CoherenceClusterSpec{Roles: []coherence.CoherenceRoleSpec{roleOne, roleTwo}},
	}
	g.Expect(cluster.GetClusterSize()).To(Equal(2))
}

func TestTwoRoleWithDifferentReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	roleOne := coherence.CoherenceRoleSpec{Role: "foo"}
	roleOne.SetReplicas(10)
	roleTwo := coherence.CoherenceRoleSpec{Role: "bar"}
	roleTwo.SetReplicas(20)
	cluster := coherence.CoherenceCluster{
		Spec: coherence.CoherenceClusterSpec{Roles: []coherence.CoherenceRoleSpec{roleOne, roleTwo}},
	}
	g.Expect(cluster.GetClusterSize()).To(Equal(30))
}

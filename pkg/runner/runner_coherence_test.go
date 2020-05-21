/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCoherenceClusterName(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Cluster: pointer.StringPtr("test-cluster"),
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.cluster="),
		"-Dcoherence.cluster=test-cluster")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceCacheConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				CacheConfig: pointer.StringPtr("test-config.xml"),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.cacheconfig="),
		"-Dcoherence.cacheconfig=test-config.xml")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceOperationalConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				OverrideConfig: pointer.StringPtr("test-override.xml"),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.k8s.override="),
		"-Dcoherence.k8s.override=test-override.xml")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceStorageEnabledTrue(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				StorageEnabled: pointer.BoolPtr(true),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=true")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceStorageEnabledFalse(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				StorageEnabled: pointer.BoolPtr(false),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=false")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceExcludeFromWKATrue(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				ExcludeFromWKA: pointer.BoolPtr(true),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceLogLevel(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				LogLevel: pointer.Int32Ptr(9),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.log.level="),
		"-Dcoherence.log.level=9")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

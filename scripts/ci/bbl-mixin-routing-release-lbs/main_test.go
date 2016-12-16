package main_test

import (
	"io/ioutil"
	"os/exec"

	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	expectedYaml = `azs:
- name: z1
  cloud_properties:
    availability_zone: us-east-1a
- name: z2
  cloud_properties:
    availability_zone: us-east-1c
- name: z3
  cloud_properties:
    availability_zone: us-east-1d
- name: z4
  cloud_properties:
    availability_zone: us-east-1e
vm_types:
- name: m3.medium
  cloud_properties:
    instance_type: m3.medium
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m3.large
  cloud_properties:
    instance_type: m3.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m3.xlarge
  cloud_properties:
    instance_type: m3.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m3.2xlarge
  cloud_properties:
    instance_type: m3.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.large
  cloud_properties:
    instance_type: m4.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.xlarge
  cloud_properties:
    instance_type: m4.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.2xlarge
  cloud_properties:
    instance_type: m4.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.4xlarge
  cloud_properties:
    instance_type: m4.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.10xlarge
  cloud_properties:
    instance_type: m4.10xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.large
  cloud_properties:
    instance_type: c3.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.xlarge
  cloud_properties:
    instance_type: c3.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.2xlarge
  cloud_properties:
    instance_type: c3.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.4xlarge
  cloud_properties:
    instance_type: c3.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.8xlarge
  cloud_properties:
    instance_type: c3.8xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.large
  cloud_properties:
    instance_type: c4.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.xlarge
  cloud_properties:
    instance_type: c4.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.2xlarge
  cloud_properties:
    instance_type: c4.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.4xlarge
  cloud_properties:
    instance_type: c4.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.8xlarge
  cloud_properties:
    instance_type: c4.8xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.large
  cloud_properties:
    instance_type: r3.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.xlarge
  cloud_properties:
    instance_type: r3.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.2xlarge
  cloud_properties:
    instance_type: r3.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.4xlarge
  cloud_properties:
    instance_type: r3.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.8xlarge
  cloud_properties:
    instance_type: r3.8xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.nano
  cloud_properties:
    instance_type: t2.nano
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.micro
  cloud_properties:
    instance_type: t2.micro
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.small
  cloud_properties:
    instance_type: t2.small
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.medium
  cloud_properties:
    instance_type: t2.medium
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.large
  cloud_properties:
    instance_type: t2.large
    ephemeral_disk:
      size: 1024
      type: gp2
disk_types:
- name: 1GB
  disk_size: 1024
  cloud_properties:
    type: gp2
    encrypted: true
- name: 5GB
  disk_size: 5120
  cloud_properties:
    type: gp2
    encrypted: true
- name: 10GB
  disk_size: 10240
  cloud_properties:
    type: gp2
    encrypted: true
- name: 50GB
  disk_size: 51200
  cloud_properties:
    type: gp2
    encrypted: true
- name: 100GB
  disk_size: 102400
  cloud_properties:
    type: gp2
    encrypted: true
- name: 500GB
  disk_size: 512000
  cloud_properties:
    type: gp2
    encrypted: true
- name: 1TB
  disk_size: 1048576
  cloud_properties:
    type: gp2
    encrypted: true
compilation:
  workers: 6
  network: private
  az: z1
  reuse_compilation_vms: true
  vm_type: c3.large
  vm_extensions:
  - 100GB_ephemeral_disk
networks:
- name: private
  type: manual
  subnets:
  - az: z1
    gateway: 10.0.16.1
    range: 10.0.16.0/20
    reserved:
    - 10.0.16.2-10.0.16.3
    - 10.0.31.255
    static:
    - 10.0.31.190-10.0.31.254
    cloud_properties:
      subnet: subnet-51e56e0a
      security_groups:
      - sg-832e75fe
  - az: z2
    gateway: 10.0.32.1
    range: 10.0.32.0/20
    reserved:
    - 10.0.32.2-10.0.32.3
    - 10.0.47.255
    static:
    - 10.0.47.190-10.0.47.254
    cloud_properties:
      subnet: subnet-4d188e60
      security_groups:
      - sg-832e75fe
  - az: z3
    gateway: 10.0.48.1
    range: 10.0.48.0/20
    reserved:
    - 10.0.48.2-10.0.48.3
    - 10.0.63.255
    static:
    - 10.0.63.190-10.0.63.254
    cloud_properties:
      subnet: subnet-2cc38265
      security_groups:
      - sg-832e75fe
  - az: z4
    gateway: 10.0.64.1
    range: 10.0.64.0/20
    reserved:
    - 10.0.64.2-10.0.64.3
    - 10.0.79.255
    static:
    - 10.0.79.190-10.0.79.254
    cloud_properties:
      subnet: subnet-4a399576
      security_groups:
      - sg-832e75fe
vm_extensions:
- name: 5GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 5120
      type: gp2
- name: 10GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 10240
      type: gp2
- name: 50GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 51200
      type: gp2
- name: 100GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 102400
      type: gp2
- name: 500GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 512000
      type: gp2
- name: 1TB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 1048576
      type: gp2
- name: router-lb
  cloud_properties:
    elbs:
    - stack-bbl-CFRouter-RN1H9557ENFE
    security_groups:
    - sg-475b003a
    - sg-832e75fe
- name: ssh-proxy-lb
  cloud_properties:
    elbs:
    - stack-bbl-CFSSHPro-12PCMYJUPS8HK
    security_groups:
    - sg-555b0028
    - sg-832e75fe
- name: tcp-router-lb
  cloud_properties:
    elbs:
    - some-tcp-router-lb
    security_groups:
    - some-tcp-router-security-group
    - some-internal-security-group`

	inputYaml = `azs:
- name: z1
  cloud_properties:
    availability_zone: us-east-1a
- name: z2
  cloud_properties:
    availability_zone: us-east-1c
- name: z3
  cloud_properties:
    availability_zone: us-east-1d
- name: z4
  cloud_properties:
    availability_zone: us-east-1e
vm_types:
- name: m3.medium
  cloud_properties:
    instance_type: m3.medium
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m3.large
  cloud_properties:
    instance_type: m3.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m3.xlarge
  cloud_properties:
    instance_type: m3.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m3.2xlarge
  cloud_properties:
    instance_type: m3.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.large
  cloud_properties:
    instance_type: m4.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.xlarge
  cloud_properties:
    instance_type: m4.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.2xlarge
  cloud_properties:
    instance_type: m4.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.4xlarge
  cloud_properties:
    instance_type: m4.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: m4.10xlarge
  cloud_properties:
    instance_type: m4.10xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.large
  cloud_properties:
    instance_type: c3.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.xlarge
  cloud_properties:
    instance_type: c3.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.2xlarge
  cloud_properties:
    instance_type: c3.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.4xlarge
  cloud_properties:
    instance_type: c3.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c3.8xlarge
  cloud_properties:
    instance_type: c3.8xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.large
  cloud_properties:
    instance_type: c4.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.xlarge
  cloud_properties:
    instance_type: c4.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.2xlarge
  cloud_properties:
    instance_type: c4.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.4xlarge
  cloud_properties:
    instance_type: c4.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: c4.8xlarge
  cloud_properties:
    instance_type: c4.8xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.large
  cloud_properties:
    instance_type: r3.large
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.xlarge
  cloud_properties:
    instance_type: r3.xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.2xlarge
  cloud_properties:
    instance_type: r3.2xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.4xlarge
  cloud_properties:
    instance_type: r3.4xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: r3.8xlarge
  cloud_properties:
    instance_type: r3.8xlarge
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.nano
  cloud_properties:
    instance_type: t2.nano
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.micro
  cloud_properties:
    instance_type: t2.micro
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.small
  cloud_properties:
    instance_type: t2.small
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.medium
  cloud_properties:
    instance_type: t2.medium
    ephemeral_disk:
      size: 1024
      type: gp2
- name: t2.large
  cloud_properties:
    instance_type: t2.large
    ephemeral_disk:
      size: 1024
      type: gp2
disk_types:
- name: 1GB
  disk_size: 1024
  cloud_properties:
    type: gp2
    encrypted: true
- name: 5GB
  disk_size: 5120
  cloud_properties:
    type: gp2
    encrypted: true
- name: 10GB
  disk_size: 10240
  cloud_properties:
    type: gp2
    encrypted: true
- name: 50GB
  disk_size: 51200
  cloud_properties:
    type: gp2
    encrypted: true
- name: 100GB
  disk_size: 102400
  cloud_properties:
    type: gp2
    encrypted: true
- name: 500GB
  disk_size: 512000
  cloud_properties:
    type: gp2
    encrypted: true
- name: 1TB
  disk_size: 1048576
  cloud_properties:
    type: gp2
    encrypted: true
compilation:
  workers: 6
  network: private
  az: z1
  reuse_compilation_vms: true
  vm_type: c3.large
  vm_extensions:
  - 100GB_ephemeral_disk
networks:
- name: private
  type: manual
  subnets:
  - az: z1
    gateway: 10.0.16.1
    range: 10.0.16.0/20
    reserved:
    - 10.0.16.2-10.0.16.3
    - 10.0.31.255
    static:
    - 10.0.31.190-10.0.31.254
    cloud_properties:
      subnet: subnet-51e56e0a
      security_groups:
      - sg-832e75fe
  - az: z2
    gateway: 10.0.32.1
    range: 10.0.32.0/20
    reserved:
    - 10.0.32.2-10.0.32.3
    - 10.0.47.255
    static:
    - 10.0.47.190-10.0.47.254
    cloud_properties:
      subnet: subnet-4d188e60
      security_groups:
      - sg-832e75fe
  - az: z3
    gateway: 10.0.48.1
    range: 10.0.48.0/20
    reserved:
    - 10.0.48.2-10.0.48.3
    - 10.0.63.255
    static:
    - 10.0.63.190-10.0.63.254
    cloud_properties:
      subnet: subnet-2cc38265
      security_groups:
      - sg-832e75fe
  - az: z4
    gateway: 10.0.64.1
    range: 10.0.64.0/20
    reserved:
    - 10.0.64.2-10.0.64.3
    - 10.0.79.255
    static:
    - 10.0.79.190-10.0.79.254
    cloud_properties:
      subnet: subnet-4a399576
      security_groups:
      - sg-832e75fe
vm_extensions:
- name: 5GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 5120
      type: gp2
- name: 10GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 10240
      type: gp2
- name: 50GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 51200
      type: gp2
- name: 100GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 102400
      type: gp2
- name: 500GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 512000
      type: gp2
- name: 1TB_ephemeral_disk
  cloud_properties:
    ephemeral_disk:
      size: 1048576
      type: gp2
- name: router-lb
  cloud_properties:
    elbs:
    - stack-bbl-CFRouter-RN1H9557ENFE
    security_groups:
    - sg-475b003a
    - sg-832e75fe
- name: ssh-proxy-lb
  cloud_properties:
    elbs:
    - stack-bbl-CFSSHPro-12PCMYJUPS8HK
    security_groups:
    - sg-555b0028
    - sg-832e75fe`
)

var _ = Describe("main", func() {
	var (
		pathToMixer string
	)

	BeforeEach(func() {
		var err error
		pathToMixer, err = gexec.Build("github.com/cloudfoundry/runtime-ci/scripts/ci/bbl-mixin-routing-release-lbs")
		Expect(err).NotTo(HaveOccurred())
	})

	It("it mixes the tcp-router vm extension into the given cloud-config", func() {
		f, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())

		_, err = f.Write([]byte(inputYaml))
		Expect(err).NotTo(HaveOccurred())

		f.Close()

		args := []string{
			"--cloud-config", f.Name(),
			"--tcp-router-elb-id", "some-tcp-router-lb",
			"--tcp-router-security-group-id", "some-tcp-router-security-group",
			"--internal-security-group-id", "some-internal-security-group",
		}

		out, err := exec.Command(pathToMixer, args...).Output()
		Expect(err).NotTo(HaveOccurred())

		Expect(out).To(gomegamatchers.MatchYAML(expectedYaml))
	})

	It("it overrides the tcp-router-lb vm extension if it already exists", func() {
		f, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())

		_, err = f.Write([]byte(`---
vm_extensions:
- name: tcp-router-lb
  cloud_properties:
    elbs:
      - some-original-lb-id
    security_groups:
      - some-original-router-sg-id
      - some-original-internal-sg-id`))
		Expect(err).NotTo(HaveOccurred())

		f.Close()

		args := []string{
			"--cloud-config", f.Name(),
			"--tcp-router-elb-id", "some-new-router-lb",
			"--tcp-router-security-group-id", "some-new-router-sg-id",
			"--internal-security-group-id", "some-new-internal-sg-id",
		}

		out, err := exec.Command(pathToMixer, args...).CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		Expect(out).To(gomegamatchers.MatchYAML([]byte(`---
vm_extensions:
  - name: tcp-router-lb
    cloud_properties:
      elbs:
        - some-new-router-lb
      security_groups:
        - some-new-router-sg-id
        - some-new-internal-sg-id`)))
	})

	Context("failure cases", func() {
		It("fails when the cloud config is not a file", func() {
			args := []string{
				"--cloud-config", "/some/bogus/path",
				"--tcp-router-elb-id", "some-new-router-lb",
				"--tcp-router-security-group-id", "some-new-router-sg-id",
				"--internal-security-group-id", "some-new-internal-sg-id",
			}

			session, err := gexec.Start(exec.Command(pathToMixer, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())

			Expect(session.Out.Contents()).To(ContainSubstring("no such file or directory"))
			Expect(session.Out.Contents()).NotTo(ContainSubstring("panic"))
			Expect(session.ExitCode()).To(Equal(1))
		})

		It("fails when the cloud config is not a valid YAML file", func() {
			f, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = f.Write([]byte("%%%%%%%"))
			Expect(err).NotTo(HaveOccurred())

			args := []string{
				"--cloud-config", f.Name(),
				"--tcp-router-elb-id", "some-new-router-lb",
				"--tcp-router-security-group-id", "some-new-router-sg-id",
				"--internal-security-group-id", "some-new-internal-sg-id",
			}

			session, err := gexec.Start(exec.Command(pathToMixer, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())

			Expect(session.Out.Contents()).To(ContainSubstring(`error parsing cloud-config: "yaml: could not find expected directive name"`))
			Expect(session.Out.Contents()).NotTo(ContainSubstring("panic"))
			Expect(session.ExitCode()).To(Equal(1))
		})
	})
})

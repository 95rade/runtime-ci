package commandgenerator_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator/commandgeneratorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commandgenerator", func() {
	var nodes int
	var env *commandgeneratorfakes.FakeEnvironment

	BeforeEach(func() {
		env = &commandgeneratorfakes.FakeEnvironment{}
		nodes = 10
		env.GetNodesReturns(nodes, nil)
		env.GetSkipBackendCompatibilityReturns("backend_compatibility", nil)
		env.GetSkipDiegoDockerReturns("docker", nil)
		env.GetSkipInternetDependentReturns("internet_dependent", nil)
		env.GetSkipLoggingReturns("logging", nil)
		env.GetSkipOperatorReturns("operator", nil)
		env.GetSkipRouteServicesReturns("route_services", nil)
		env.GetSkipSecurityGroupsReturns("security_groups", nil)
		env.GetSkipServicesReturns("services", nil)
		env.GetSkipDiegoSSHReturns("ssh", nil)
		env.GetSkipV3Returns("v3", nil)
	})

	Context("When the path to CATs is set", func() {
		BeforeEach(func() {
			env.GetCatsPathReturns(".")
		})

		It("Should generate a command to run CATS", func() {
			cmd, args, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("bin/test"))

			Expect(strings.Join(args, " ")).To(Equal(
				fmt.Sprintf("-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes=%d -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing", nodes),
			))

			env.GetCatsPathReturns("/path/to/cats")
			cmd, _, err = commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/path/to/cats/bin/test"))
		})

		Context("when the env returns an error fetching the number of nodes", func() {
			var expectedError error
			BeforeEach(func() {
				expectedError = fmt.Errorf("some error")
				env.GetNodesReturns(0, expectedError)
			})

			It("propogates the error", func() {
				_, _, err := commandgenerator.GenerateCmd(env)
				Expect(err.Error()).To(Equal(expectedError.Error()))
			})
		})

		Context("when the node count is unset", func() {
			BeforeEach(func() {
				env.GetNodesReturns(0, nil)
			})
			It("sets the default node count", func() {
				_, args, _ := commandgenerator.GenerateCmd(env)
				Expect(args).To(ContainElement("-nodes=2"))
			})
		})

		Context("when there are optional skipPackage env vars set", func() {
			Context("to not skip", func() {
				BeforeEach(func() {
					env.GetSkipBackendCompatibilityReturns("", nil)
					env.GetSkipDiegoDockerReturns("", nil)
					env.GetSkipInternetDependentReturns("", nil)
					env.GetSkipLoggingReturns("", nil)
					env.GetSkipOperatorReturns("", nil)
					env.GetSkipRouteServicesReturns("", nil)
					env.GetSkipSecurityGroupsReturns("", nil)
					env.GetSkipServicesReturns("", nil)
					env.GetSkipDiegoSSHReturns("", nil)
					env.GetSkipV3Returns("", nil)
				})

				It("should generate a command with the correct list of skipPackage flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(ContainSubstring(" -skipPackage=helpers "))
				})
			})

			Context("to skip", func() {
				It("should generate a command with the correct list of skipPackage flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(ContainSubstring(
						" -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 ",
					))
				})
			})

			Context("when the env returns an error", func() {
				var expectedError error
				BeforeEach(func() {
					expectedError = fmt.Errorf("some error")
					env.GetSkipV3Returns("", expectedError)
				})

				It("propogates the error", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err.Error()).To(Equal(expectedError.Error()))
				})
			})
		})

		Context("when there are optional skip env vars set", func() {
			Context("to skip SSO", func() {
				BeforeEach(func() {
					env.GetSkipSSOReturns("SSO", nil)
				})

				It("generates a command that skips tests with the given tag", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).To(ContainElement(ContainSubstring("-skip=SSO|")))
				})
			})

			Context("to not skip SSO", func() {
				BeforeEach(func() {
					env.GetSkipSSOReturns("", nil)
				})

				It("generates a command that does not include the given tag in the skips", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).ToNot(ContainElement(ContainSubstring("-skip=SSO|")))
				})
			})

			Context("and the env returns an error", func() {
				var expectedError error
				BeforeEach(func() {
					expectedError = fmt.Errorf("some error")
					env.GetSkipSSOReturns("", expectedError)
				})

				It("propogates the error", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err.Error()).To(Equal(expectedError.Error()))
					Expect(env.GetSkipSSOCallCount())
				})
			})
		})

		Describe("the BACKEND parameter", func() {
			Context("is set to diego", func() {
				BeforeEach(func() {
					env.GetBackendReturns("diego", nil)
				})

				It("should generate a command that skips NO_DIEGO_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).To(ContainElement("-skip=NO_DIEGO_SUPPORT"))
				})

			})

			Context("is set to dea", func() {
				BeforeEach(func() {
					env.GetBackendReturns("dea", nil)
				})

				It("should generate a command that skips NO_DEA_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))
					Expect(args).To(ContainElement("-skip=NO_DEA_SUPPORT"))
				})
			})

			Context("isn't set", func() {
				BeforeEach(func() {
					env.GetBackendReturns("", nil)
				})
				It("should generate a command that skips both NO_DIEGO_SUPPORT and NO_DEA_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))
					Expect(args).To(ContainElement("-skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT"))
				})
			})

			Context("is not a valid value", func() {
				BeforeEach(func() {
					env.GetBackendReturns("", fmt.Errorf("some error"))
				})
				It("returns an error", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("some error"))
				})
			})
		})
	})

	Context("When the path to CATS isn't explicitly provided", func() {
		BeforeEach(func() {
			env.GetGoPathReturns("/go")
		})

		It("Should return a sane default command path for use in Concourse", func() {
			cmd, _, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/go/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"))
		})
	})

	Context("When everything is returning an error", func() {
		BeforeEach(func() {
			env.GetNodesReturns(0, fmt.Errorf("NODES"))
			env.GetSkipBackendCompatibilityReturns("", fmt.Errorf("INCLUDE_BACKEND_COMPATIBILITY"))
			env.GetSkipDiegoDockerReturns("", fmt.Errorf("INCLUDE_DIEGO_DOCKER"))
			env.GetSkipInternetDependentReturns("", fmt.Errorf("INCLUDE_INTERNET_DEPENDENT"))
			env.GetSkipLoggingReturns("", fmt.Errorf("INCLUDE_LOGGING"))
			env.GetSkipOperatorReturns("", fmt.Errorf("INCLUDE_OPERATOR"))
			env.GetSkipRouteServicesReturns("", fmt.Errorf("INCLUDE_ROUTE_SERVICES"))
			env.GetSkipSecurityGroupsReturns("", fmt.Errorf("INCLUDE_SECURITY_GROUPS"))
			env.GetSkipServicesReturns("", fmt.Errorf("INCLUDE_SERVICES"))
			env.GetSkipDiegoSSHReturns("", fmt.Errorf("INCLUDE_DIEGO_SSH"))
			env.GetSkipV3Returns("", fmt.Errorf("INCLUDE_V3"))
		})
		It("Aggregates all the errors", func() {
			_, _, err := commandgenerator.GenerateCmd(env)
			errorStr := err.Error()
			Expect(errorStr).To(
				And(
					ContainSubstring("NODES"),
					ContainSubstring("INCLUDE_BACKEND_COMPATIBILITY"),
					ContainSubstring("INCLUDE_DIEGO_DOCKER"),
					ContainSubstring("INCLUDE_INTERNET_DEPENDENT"),
					ContainSubstring("INCLUDE_LOGGING"),
					ContainSubstring("INCLUDE_OPERATOR"),
					ContainSubstring("INCLUDE_ROUTE_SERVICES"),
					ContainSubstring("INCLUDE_SECURITY_GROUPS"),
					ContainSubstring("INCLUDE_SERVICES"),
					ContainSubstring("INCLUDE_DIEGO_SSH"),
					ContainSubstring("INCLUDE_V3"),
				),
			)
		})
	})
})

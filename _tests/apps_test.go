package _tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apps", func() {
	Context("with a logged-in user", func() {
		appName := "apps-test"

		BeforeEach(func() {
			login(url, testUser, testPassword)
		})

		It("can't get app info", func() {
			output, err := execute("deis info -a %s", appName)
			Expect(err).To(HaveOccurred())
			Expect(output).To(ContainSubstring("NOT FOUND"))
		})

		It("can't get app logs", func() {
			output, err := execute("deis logs -a %s", appName)
			Expect(err).To(HaveOccurred())
			Expect(output).To(ContainSubstring("NOT FOUND"))
		})

		// TODO: this currently returns "Error: json: cannot unmarshal object into Go value of type []interface {}"
		XIt("can't run a command in the app environment", func() {
			output, err := execute("deis apps:run echo Hello, 世界")
			Expect(err).To(HaveOccurred())
			Expect(output).To(ContainSubstring("NOT FOUND"))
		})

		It("can create an app", func() {
			output, err := execute("deis apps:create %s", appName)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				ContainSubstring("Creating Application... done, created %s", appName),
				ContainSubstring("Git remote deis added"),
				ContainSubstring("remote available at ")))
			output, err = execute("deis apps:destroy --confirm=%s", appName)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				ContainSubstring("Destroying %s...", appName),
				ContainSubstring("done in "),
				ContainSubstring("Git remote deis removed")))
		})

		It("can create an app with no git remote", func() {
			output, err := execute("deis apps:create %s --no-remote", appName)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				ContainSubstring("Creating Application... done, created %s", appName),
				ContainSubstring("remote available at ")))
			Expect(output).NotTo(ContainSubstring("Git remote deis added"))
			output, err = execute("deis apps:destroy --app=%s --confirm=%s", appName, appName)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				ContainSubstring("Destroying %s...", appName),
				ContainSubstring("done in ")))
			Expect(output).NotTo(ContainSubstring("Git remote deis removed"))
		})

		It("can create an app with a custom buildpack", func() {
			output, err := execute("deis apps:create %s --buildpack https://example.com", appName)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				ContainSubstring("Creating Application... done, created %s", appName),
				ContainSubstring("Git remote deis added"),
				ContainSubstring("remote available at ")))
			output, err = execute("deis config:list")
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(ContainSubstring("BUILDPACK_URL"))
			output, err = execute("deis apps:destroy --app=%s --confirm=%s", appName, appName)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				ContainSubstring("Destroying %s...", appName),
				ContainSubstring("done in "),
				ContainSubstring("Git remote deis removed")))
		})
	})

	// Context("with a deployed app", func() {
	//
	// 	appName := "apps-test"
	// 	repository := "https://github.com/deis/example-go.git"
	//
	// TODO: can't have an Expect outside an It clause...need to refactor
	// 	output, err := execute("git clone %s", repository)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(output).To(SatisfyAll(
	// 		ContainSubstring("Cloning into "),
	// 		ContainSubstring("done.")))
	// 	// TODO: change directory to cloned app dir
	// 	output, err = execute("deis apps:create %s", appName)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(output).To(SatisfyAll(
	// 		ContainSubstring("Creating Application... done, created %s", appName),
	// 		ContainSubstring("Git remote deis added"),
	// 		ContainSubstring("remote available at ")))
	// 	output, err = execute("git push deis master")
	// 	Expect(err).NotTo(HaveOccurred())
	// 	Expect(output).To(SatisfyAll(
	// 		ContainSubstring("-----> Launching..."),
	// 		ContainSubstring("done, %s:v2 deployed to Deis", appName)))
	//
	// 	It("can't create an existing app", func() {
	// 		output, err = execute("deis apps:create %s", appName)
	// 		Expect(err).To(HaveOccurred())
	// 		Expect(output).To(ContainSubstring("This field must be unique"))
	// 	})
	//
	// 	It("can get app info", func() {
	// 		output, err := execute("deis info")
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(output).To(SatisfyAll(
	// 			HavePrefix("=== %s Application", appName),
	// 			ContainSubstring("=== %s Processes", appName),
	// 			ContainSubstring(".1 up (v"),
	// 			ContainSubstring("=== %s Domains", appName)))
	// 	})
	//
	// 	It("can get app logs", func() {
	// 		output, err := execute("deis logs")
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(output).To(SatisfyAll(
	// 			ContainSubstring("%s[deis-controller]: %s created initial release",
	// 				appName, username),
	// 			ContainSubstring("%s[deis-controller]: %s deployed", appName, username),
	// 			ContainSubstring("%s[deis-controller]: %s scaled containers",
	// 				appName, username)))
	// 	})
	//
	// 	// TODO: how to test "deis open" which spawns a browser?
	// 	XIt("can open the app's URL", func() {
	// 		_, err := execute("deis open")
	// 		Expect(err).NotTo(HaveOccurred())
	// 	})
	//
	// 	It("can't open a bogus app URL", func() {
	// 		output, err := execute("deis open -a bogus-appname")
	// 		Expect(err).To(HaveOccurred())
	// 		Expect(output).To(ContainSubstring("404 NOT FOUND"))
	// 	})
	//
	// 	It("can run a command in the app environment", func() {
	// 		output, err := execute("deis apps:run echo Hello, 世界")
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(output).To(SatisfyAll(
	// 			HavePrefix("Running 'echo Hello, 世界'..."),
	// 			HaveSuffix("Hello, 世界\n")))
	// 	})
	//
	// 	// TODO: this requires a second user account
	// 	XIt("can transfer the app to another owner", func() {
	// 	})
	// })

})

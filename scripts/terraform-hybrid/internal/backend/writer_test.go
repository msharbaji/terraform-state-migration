package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	TestWorkspaceDir = "/tmp/deploy/provider/test-module"
	TestCallerName   = "test-caller"
	LocalBackendPath = "/tmp/terraform-states"
)

func TestBackendWriter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TerraformBackendWriter Suite")
}

var _ = Describe("TerraformBackendWriter", func() {
	var (
		tbw           *TerraformBackendWriter
		workspaceDir  string
		callerName    string
		backendConfig *config.TerraformHybridConfig
	)

	// Helper function to check file content
	checkBackendFileContent := func(workspaceDir string, expectedContent []string) {
		backendFile := filepath.Join(workspaceDir, "backend.tf")
		Expect(backendFile).To(BeAnExistingFile())

		content, err := os.ReadFile(backendFile)
		Expect(err).To(BeNil())

		for _, substring := range expectedContent {
			Expect(string(content)).To(ContainSubstring(substring))
		}
	}

	BeforeEach(func() {
		tbw = &TerraformBackendWriter{}
		workspaceDir = TestWorkspaceDir
		callerName = TestCallerName

		// Ensure the directory exists before tests run
		err := os.MkdirAll(workspaceDir, 0755)
		Expect(err).To(BeNil())

		backendConfig = &config.TerraformHybridConfig{
			Global: config.GlobalConfig{
				Backend:     nil,
				BackendType: config.LocalBackendType,
			},
		}
	})

	AfterEach(func() {
		// Clean up the test directory after each test run
		err := os.RemoveAll("/tmp/deploy/provider")
		Expect(err).To(BeNil())
	})

	DescribeTable("should write backend configurations",
		func(backendType config.BackendType, backend interface{}, expectedContent []string) {
			backendConfig.Global.BackendType = backendType
			backendConfig.Global.Backend = backend
			err := tbw.WriteBackend(backendConfig, workspaceDir, callerName)
			Expect(err).To(BeNil())

			checkBackendFileContent(workspaceDir, expectedContent)
		},
		Entry("Local Backend", config.LocalBackendType, &config.LocalBackendConfig{Path: LocalBackendPath},
			[]string{"backend \"local\"", fmt.Sprintf("path = \"%s/test-module/terraform.tfstate\"", LocalBackendPath)}),
		Entry("Cloud Storage Backend", config.BackendTypeCloudStorage, &config.CloudStorageBackendConfig{
			Type: "gcs", Region: "us-east-1", BucketName: "my-terraform-state-bucket"},
			[]string{"backend \"gcs\"", "bucket  = \"my-terraform-state-bucket\""}),
		Entry("Postgres Backend", config.BackendTypePostgres, &config.PostgresBackendConfig{
			ConnectionString: "postgres://user:password@localhost:5432/terraform", SchemaName: "terraform_state"},
			[]string{"backend \"pg\"", "conn_str     = \"postgres://user:password@localhost:5432/terraform\"", "schema_name  = \"terraform_state\""}),
	)

	Context("when the backend type is unsupported", func() {
		It("should return an error", func() {
			backendConfig.Global.BackendType = config.BackendType("unsupported")
			err := tbw.WriteBackend(backendConfig, workspaceDir, callerName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported backend type"))
		})
	})

	Describe("getRelativePathUnderProvider", func() {
		It("should calculate the relative path under deploy/provider", func() {
			absPath := TestWorkspaceDir
			expectedRelativePath := "test-module"
			relativePath, err := tbw.getRelativePathUnderProvider(absPath)
			Expect(err).To(BeNil())
			Expect(relativePath).To(Equal(expectedRelativePath))
		})

		It("should return an error if the workspace is not under deploy/provider", func() {
			invalidPath := "/some/random/path"
			_, err := tbw.getRelativePathUnderProvider(invalidPath)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("workspace directory does not seem to be under 'deploy/provider'"))
		})
	})
})

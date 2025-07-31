package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
)

// TestS3BucketDataConfiguration tests the data S3 bucket configuration
func TestS3BucketDataConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	// Test terraform plan output
	t.Run("DataBucketPlanValidation", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify data bucket is in the plan
		assert.Contains(t, planOutput, "aws_s3_bucket.data")
		assert.Contains(t, planOutput, "force_destroy = true")
		
		// Verify bucket naming convention
		assert.Contains(t, planOutput, "${local.resource_prefix.value}-data")
	})

	// Test security vulnerabilities in data bucket (intentionally insecure)
	t.Run("DataBucketSecurityVulnerabilities", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Data bucket should NOT have encryption (vulnerability)
		assert.NotContains(t, planOutput, "server_side_encryption_configuration")
		
		// Data bucket should NOT have versioning (vulnerability)
		versioningCount := strings.Count(planOutput, "versioning")
		dataVersioning := strings.Contains(planOutput, "aws_s3_bucket.data") && strings.Contains(planOutput, "versioning")
		assert.False(t, dataVersioning, "Data bucket should not have versioning configured")
		
		// Data bucket should NOT have access logging (vulnerability)
		loggingCount := strings.Count(planOutput, "logging")
		dataLogging := strings.Contains(planOutput, "aws_s3_bucket.data") && strings.Contains(planOutput, "logging")
		assert.False(t, dataLogging, "Data bucket should not have access logging configured")
	})

	// Test tags structure for data bucket
	t.Run("DataBucketTagsValidation", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify git-related tags are present
		assert.Contains(t, planOutput, "git_commit")
		assert.Contains(t, planOutput, "git_file")
		assert.Contains(t, planOutput, "git_org")
		assert.Contains(t, planOutput, "yor_trace")
		
		// Verify merge pattern is used for tags
		assert.Contains(t, planOutput, "tags = merge")
		
		// Verify specific tag values
		assert.Contains(t, planOutput, "git_org              = \"bridgecrewio\"")
		assert.Contains(t, planOutput, "git_repo             = \"terragoat\"")
	})
}

// TestS3BucketWarehouseConfiguration tests warehouse bucket (duplicate name issue)
func TestS3BucketWarehouseConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("WarehouseBucketNamingConflict", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Both data and warehouse buckets use the same name - this is a bug
		assert.Contains(t, planOutput, "aws_s3_bucket.warehouse")
		assert.Contains(t, planOutput, "aws_s3_bucket.data")
		
		// Both use the same bucket name pattern - potential conflict
		dataNameOccurrences := strings.Count(planOutput, "${local.resource_prefix.value}-data")
		assert.GreaterOrEqual(t, dataNameOccurrences, 2, "Both data and warehouse buckets use same name pattern")
	})

	t.Run("WarehouseBucketBasicConfig", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify basic warehouse bucket configuration
		assert.Contains(t, planOutput, "force_destroy = true")
		assert.Contains(t, planOutput, "git_commit")
		assert.Contains(t, planOutput, "yor_trace")
	})
}

// TestS3BucketObjectConfiguration tests S3 object configuration
func TestS3BucketObjectConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test", 
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("S3ObjectConfiguration", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify S3 object configuration
		assert.Contains(t, planOutput, "aws_s3_bucket_object.data_object")
		assert.Contains(t, planOutput, "customer-master.xlsx")
		assert.Contains(t, planOutput, "resources/customer-master.xlsx")
		
		// Verify dependency on data bucket
		assert.Contains(t, planOutput, "bucket = aws_s3_bucket.data.id")
	})

	t.Run("S3ObjectTags", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify object tags
		assert.Contains(t, planOutput, "customer-master")
		assert.Contains(t, planOutput, "Environment")
		assert.Contains(t, planOutput, "yor_trace")
	})
}

// TestS3BucketFinancialsConfiguration tests financials bucket
func TestS3BucketFinancialsConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("FinancialsBucketACL", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify private ACL is configured
		assert.Contains(t, planOutput, "aws_s3_bucket.financials")
		assert.Contains(t, planOutput, "acl           = \"private\"")
	})

	t.Run("FinancialsBucketSecurityGaps", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Financials bucket should NOT have encryption (vulnerability)
		financialsEncryption := strings.Contains(planOutput, "aws_s3_bucket.financials") && 
			strings.Contains(planOutput, "server_side_encryption_configuration")
		assert.False(t, financialsEncryption, "Financials bucket should not have encryption")
		
		// Financials bucket should NOT have versioning (vulnerability)
		financialsVersioning := strings.Contains(planOutput, "aws_s3_bucket.financials") && 
			strings.Contains(planOutput, "versioning")
		assert.False(t, financialsVersioning, "Financials bucket should not have versioning")
	})

	t.Run("FinancialsBucketNaming", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify unique naming for financials bucket
		assert.Contains(t, planOutput, "${local.resource_prefix.value}-financials")
	})
}

// TestS3BucketOperationsConfiguration tests operations bucket
func TestS3BucketOperationsConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("OperationsBucketVersioning", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Operations bucket SHOULD have versioning enabled
		assert.Contains(t, planOutput, "aws_s3_bucket.operations")
		assert.Contains(t, planOutput, "versioning")
		assert.Contains(t, planOutput, "enabled = true")
	})

	t.Run("OperationsBucketSecurityPartial", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify private ACL
		assert.Contains(t, planOutput, "acl    = \"private\"")
		
		// Operations bucket should NOT have encryption (still a vulnerability)
		operationsEncryption := strings.Contains(planOutput, "aws_s3_bucket.operations") && 
			strings.Contains(planOutput, "server_side_encryption_configuration")
		assert.False(t, operationsEncryption, "Operations bucket should not have encryption")
		
		// Operations bucket should NOT have access logging (vulnerability)
		operationsLogging := strings.Contains(planOutput, "aws_s3_bucket.operations") && 
			strings.Contains(planOutput, "logging")
		assert.False(t, operationsLogging, "Operations bucket should not have access logging")
	})
}

// TestS3BucketDataScienceConfiguration tests data science bucket
func TestS3BucketDataScienceConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("DataScienceBucketVersioningAndLogging", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Data science bucket SHOULD have versioning
		assert.Contains(t, planOutput, "aws_s3_bucket.data_science")
		assert.Contains(t, planOutput, "versioning")
		assert.Contains(t, planOutput, "enabled = true")
		
		// Data science bucket SHOULD have access logging
		assert.Contains(t, planOutput, "logging")
		assert.Contains(t, planOutput, "target_bucket = \"${aws_s3_bucket.logs.id}\"")
		assert.Contains(t, planOutput, "target_prefix = \"log/\"")
	})

	t.Run("DataScienceBucketEncryptionGap", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Data science bucket should NOT have encryption (still a vulnerability)
		dataScienceEncryption := strings.Contains(planOutput, "aws_s3_bucket.data_science") && 
			strings.Contains(planOutput, "server_side_encryption_configuration")
		assert.False(t, dataScienceEncryption, "Data science bucket should not have encryption")
	})

	t.Run("DataScienceBucketTagsPattern", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Data science bucket uses direct tags (not merge pattern)
		assert.Contains(t, planOutput, "git_commit")
		assert.Contains(t, planOutput, "yor_trace")
		
		// Should NOT use merge pattern like other buckets
		dataScienceMerge := strings.Contains(planOutput, "aws_s3_bucket.data_science") && 
			strings.Contains(planOutput, "tags = merge")
		assert.False(t, dataScienceMerge, "Data science bucket should use direct tags, not merge")
	})
}

// TestS3BucketLogsConfiguration tests logs bucket (most secure)
func TestS3BucketLogsConfiguration(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("LogsBucketFullSecurity", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Logs bucket SHOULD have encryption (most secure bucket)
		assert.Contains(t, planOutput, "server_side_encryption_configuration")
		assert.Contains(t, planOutput, "sse_algorithm     = \"aws:kms\"")
		assert.Contains(t, planOutput, "kms_master_key_id = \"${aws_kms_key.logs_key.arn}\"")
		
		// Logs bucket SHOULD have versioning
		assert.Contains(t, planOutput, "versioning")
		assert.Contains(t, planOutput, "enabled = true")
	})

	t.Run("LogsBucketSpecialACL", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Logs bucket has special ACL for log delivery
		assert.Contains(t, planOutput, "acl    = \"log-delivery-write\"")
	})

	t.Run("LogsBucketNamingAndTags", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify naming convention
		assert.Contains(t, planOutput, "${local.resource_prefix.value}-logs")
		
		// Verify merge tags pattern
		assert.Contains(t, planOutput, "tags = merge")
		assert.Contains(t, planOutput, "Name        = \"${local.resource_prefix.value}-logs\"")
		assert.Contains(t, planOutput, "Environment = local.resource_prefix.value")
	})
}

// TestS3SecurityComplianceMatrix tests security across all buckets
func TestS3SecurityComplianceMatrix(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("EncryptionComplianceAudit", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Only logs bucket should have encryption
		encryptionOccurrences := strings.Count(planOutput, "server_side_encryption_configuration")
		assert.Equal(t, 1, encryptionOccurrences, "Only logs bucket should have encryption")
		
		// Verify logs bucket has KMS encryption
		assert.Contains(t, planOutput, "aws:kms")
	})

	t.Run("VersioningComplianceAudit", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Count versioning occurrences - should be 3 (operations, data_science, logs)
		versioningOccurrences := strings.Count(planOutput, "versioning")
		assert.Equal(t, 3, versioningOccurrences, "Three buckets should have versioning enabled")
	})

	t.Run("AccessLoggingComplianceAudit", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Only data_science bucket should have access logging
		loggingOccurrences := strings.Count(planOutput, "logging")
		assert.Equal(t, 1, loggingOccurrences, "Only data_science bucket should have access logging")
		
		// Verify logging target is logs bucket
		assert.Contains(t, planOutput, "target_bucket = \"${aws_s3_bucket.logs.id}\"")
	})
}

// TestS3ResourceDependencies tests inter-resource dependencies
func TestS3ResourceDependencies(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("S3ObjectBucketDependency", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// S3 object depends on data bucket
		assert.Contains(t, planOutput, "bucket = aws_s3_bucket.data.id")
	})

	t.Run("LoggingBucketDependency", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Data science bucket logging depends on logs bucket
		assert.Contains(t, planOutput, "target_bucket = \"${aws_s3_bucket.logs.id}\"")
	})

	t.Run("KMSKeyDependency", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Logs bucket encryption depends on KMS key
		assert.Contains(t, planOutput, "kms_master_key_id = \"${aws_kms_key.logs_key.arn}\"")
	})
}

// TestS3EdgeCases tests edge cases and error conditions
func TestS3EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("EmptyResourcePrefix", func(t *testing.T) {
		terraformOptions := &terraform.Options{
			TerraformDir: "../..",
			Vars: map[string]interface{}{
				"resource_prefix": "",
			},
			NoColor: true,
		}
		
		// Should handle empty prefix (might cause naming issues)
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		assert.Contains(t, planOutput, "aws_s3_bucket")
		
		// Verify bucket names would be just suffixes
		assert.Contains(t, planOutput, "-data")
		assert.Contains(t, planOutput, "-logs")
	})

	t.Run("LongResourcePrefix", func(t *testing.T) {
		longPrefix := "very-long-resource-prefix-that-might-exceed-aws-bucket-naming-limits-and-cause-issues"
		terraformOptions := &terraform.Options{
			TerraformDir: "../..",
			Vars: map[string]interface{}{
				"resource_prefix": longPrefix,
			},
			NoColor: true,
		}
		
		// Should handle long prefix (might exceed AWS limits)
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		assert.Contains(t, planOutput, "aws_s3_bucket")
	})

	t.Run("SpecialCharactersInPrefix", func(t *testing.T) {
		terraformOptions := &terraform.Options{
			TerraformDir: "../..",
			Vars: map[string]interface{}{
				"resource_prefix": "test_123-special.chars",
			},
			NoColor: true,
		}
		
		// Should handle special characters in prefix
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		assert.Contains(t, planOutput, "aws_s3_bucket")
	})
}

// TestS3TagConsistencyAndStandards tests tag consistency
func TestS3TagConsistencyAndStandards(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("GitTagsConsistency", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// All resources should have consistent git tags
		requiredGitTags := []string{
			"git_commit",
			"git_file",
			"git_org",
			"git_repo",
			"yor_trace",
		}
		
		for _, tag := range requiredGitTags {
			tagOccurrences := strings.Count(planOutput, tag)
			assert.GreaterOrEqual(t, tagOccurrences, 6, "Tag %s should appear in all resources", tag)
		}
		
		// Verify consistent git values
		assert.Contains(t, planOutput, "git_org              = \"bridgecrewio\"")
		assert.Contains(t, planOutput, "git_repo             = \"terragoat\"")
		assert.Contains(t, planOutput, "git_file             = \"terraform/aws/s3.tf\"")
	})

	t.Run("TagPatternConsistency", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Most buckets use merge pattern for tags
		mergeOccurrences := strings.Count(planOutput, "tags = merge")
		assert.Equal(t, 5, mergeOccurrences, "Five buckets should use merge pattern for tags")
		
		// Data science bucket uses direct tags pattern
		assert.Contains(t, planOutput, "aws_s3_bucket.data_science")
	})
}

// TestS3PlanOutputStructure tests overall plan structure
func TestS3PlanOutputStructure(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../..",
		Vars: map[string]interface{}{
			"resource_prefix": "test",
		},
		NoColor: true,
	}

	defer terraform.Destroy(t, terraformOptions)

	t.Run("AllResourcesPresent", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Verify all expected resources are in the plan
		expectedResources := []string{
			"aws_s3_bucket.data",
			"aws_s3_bucket.warehouse",
			"aws_s3_bucket_object.data_object",
			"aws_s3_bucket.financials",
			"aws_s3_bucket.operations",
			"aws_s3_bucket.data_science",
			"aws_s3_bucket.logs",
		}
		
		for _, resource := range expectedResources {
			assert.Contains(t, planOutput, resource, "Resource %s should be present in plan", resource)
		}
	})

	t.Run("ResourceCountValidation", func(t *testing.T) {
		planOutput := terraform.InitAndPlan(t, terraformOptions)
		
		// Count S3 bucket resources (should be 6)
		bucketCount := strings.Count(planOutput, "resource \"aws_s3_bucket\"")
		assert.Equal(t, 6, bucketCount, "Should have exactly 6 S3 buckets")
		
		// Count S3 bucket objects (should be 1)
		objectCount := strings.Count(planOutput, "resource \"aws_s3_bucket_object\"")
		assert.Equal(t, 1, objectCount, "Should have exactly 1 S3 bucket object")
	})
}
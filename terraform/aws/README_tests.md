# S3 Terraform Configuration Tests

This test suite provides comprehensive testing for the S3 Terraform configuration using the Terratest framework.

## Overview

The test suite validates:
- S3 bucket configurations and security settings
- Resource dependencies and naming conventions
- Tag consistency and compliance
- Security vulnerabilities and compliance gaps
- Edge cases and error conditions

## Testing Framework

- **Framework**: Terratest (Go-based Terraform testing framework)
- **Assertions**: testify/assert and testify/require
- **Parallel Execution**: Tests run in parallel for faster execution

## Test Categories

### 1. Individual Bucket Tests
- `TestS3BucketDataConfiguration`: Tests the intentionally insecure data bucket
- `TestS3BucketWarehouseConfiguration`: Tests warehouse bucket (duplicate naming issue)
- `TestS3BucketObjectConfiguration`: Tests S3 object configuration
- `TestS3BucketFinancialsConfiguration`: Tests financials bucket
- `TestS3BucketOperationsConfiguration`: Tests operations bucket (partial security)
- `TestS3BucketDataScienceConfiguration`: Tests data science bucket (versioning + logging)
- `TestS3BucketLogsConfiguration`: Tests the most secure logs bucket

### 2. Security Compliance Tests
- `TestS3SecurityComplianceMatrix`: Audits security configurations across all buckets
- Validates encryption, versioning, and access logging compliance

### 3. Dependency Tests
- `TestS3ResourceDependencies`: Tests inter-resource dependencies
- Validates S3 object → bucket, logging → logs bucket, encryption → KMS key relationships

### 4. Edge Case Tests
- `TestS3EdgeCases`: Tests error conditions and edge cases
- Empty prefixes, long prefixes, special characters

### 5. Consistency Tests
- `TestS3TagConsistencyAndStandards`: Tests tag consistency across resources
- `TestS3PlanOutputStructure`: Tests overall Terraform plan structure

## Security Vulnerabilities Detected

The tests intentionally validate several security vulnerabilities present in the configuration:

### Unencrypted Buckets (5/6 buckets)
- `data` bucket: No encryption
- `warehouse` bucket: No encryption  
- `financials` bucket: No encryption
- `operations` bucket: No encryption
- `data_science` bucket: No encryption
- ✅ `logs` bucket: KMS encryption enabled

### Missing Versioning (3/6 buckets)
- `data` bucket: No versioning
- `warehouse` bucket: No versioning
- `financials` bucket: No versioning
- ✅ `operations` bucket: Versioning enabled
- ✅ `data_science` bucket: Versioning enabled
- ✅ `logs` bucket: Versioning enabled

### Missing Access Logging (5/6 buckets)
- `data` bucket: No access logging
- `warehouse` bucket: No access logging
- `financials` bucket: No access logging
- `operations` bucket: No access logging
- ✅ `data_science` bucket: Access logging to logs bucket
- `logs` bucket: No access logging (expected for log destination)

### Configuration Issues
- **Duplicate Naming**: Both `data` and `warehouse` buckets use same name pattern
- **Inconsistent Tagging**: `data_science` bucket uses direct tags vs merge pattern

## Running Tests

```bash
# Run all tests
make test

# Run tests in parallel
make test-parallel

# Clean test cache
make clean

# Run specific test
go test -v -run TestS3BucketDataConfiguration

# Run with verbose output
go test -v -timeout 30m
```

## Prerequisites

1. Go 1.19 or later
2. Terraform installed
3. AWS credentials configured (for actual deployment tests)
4. Dependencies: `go mod tidy`

## Test Structure

Each test follows this pattern:
1. **Setup**: Configure Terraform options with test variables
2. **Defer Cleanup**: Ensure resources are destroyed after test
3. **Plan Validation**: Use `terraform.InitAndPlan()` to validate configuration
4. **Assertions**: Test specific aspects of the Terraform plan output

## Notes

- Tests use `terraform plan` output validation (static analysis)
- No actual AWS resources are created during plan-only tests
- For integration tests, set appropriate AWS credentials and remove plan-only mode
- Tests run in parallel by default for faster execution
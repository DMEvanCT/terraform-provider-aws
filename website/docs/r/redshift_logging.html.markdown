---
subcategory: "Redshift"
layout: "aws"
page_title: "AWS: aws_redshift_logging"
description: |-
  Terraform resource for managing an AWS Redshift Logging configuration.
---
# Resource: aws_redshift_logging

Terraform resource for managing an AWS Redshift Logging configuration.

~> In order to prevent persistent differences when using this resource, the parent cluster should include an `ignore_changes` lifecycle configuration containing the `logging` argument. This will be necessary only until the deprecated argument is removed in a future major version.

## Example Usage

### Basic Usage

```terraform
resource "aws_redshift_cluster" "example" {
  ### other configuration ###

  # An ignore_changes lifecycle configuration is required until the deprecated
  # logging argument is removed in a future major version.
  lifecycle {
    ignore_changes = [logging]
  }
}

resource "aws_redshift_logging" "example" {
  cluster_identifier   = aws_redshift_cluster.example.id
  log_destination_type = "cloudwatch"
}
```

### S3 Destination Type

```terraform
resource "aws_redshift_logging" "example" {
  cluster_identifier   = aws_redshift_cluster.example.id
  bucket_name          = aws_s3_bucket.example.id
  log_destination_type = "s3"
  log_exports          = ["connectionlog", "userlog"]
  s3_key_prefix        = "example-prefix"
}
```

## Argument Reference

The following arguments are required:

* `cluster_identifier` - (Required) Identifier of the source cluster.

The following arguments are optional:

* `bucket_name` - (Optional) Name of an existing S3 bucket where the log files are to be stored.
* `log_destination_type` - (Optional) Log destination type. Valid values are `cloudwatch` and `s3`.
* `log_exports` - (Optional) An array of exported log types. Valid values are `connectionlog`, `useractivitylog`, and `userlog`.
* `s3_key_prefix` - (Optional) Prefix applied to the log file names.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - Identifier of the source cluster.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Redshift Logging using the `id`. For example:

```terraform
import {
  to = aws_redshift_logging.example
  id = "cluster-id-12345678"
}
```

Using `terraform import`, import Redshift Logging using the `id`. For example:

```console
% terraform import aws_redshift_logging.example cluster-id-12345678
```

# Sift

Sift provides policy enforcement fromework for application and organizations
running at scale in the cloud.

# Philosophy

Running large application in the cloud is difficult. Enterprises will in a very short
amount of time ind that they have dozens of different accounts, across multiple
providers, service many different projects.

You want the teams running these projects to have control of their cloud resources,
however with the proper guardrails in place to ensure they are not exposing their
products to risk in the cloud.

Sift takes the approach that the only way to keep up to date with ensure
proper configuration and not needlessly hindering application development, is to
create policies which provide broad boundaries around an application, with the
ability to be fosed in very narrowly when necessary.

# Tenants

* Provide reasonable boundaries for what an application can do.
* Allows for a fine grained level of control.
* Supports exceptions and custom policies at a granular level.
* Build to scale across multiple providers, with multiple projects.
* Optomized for speed, caches and parallizes when possible.
* Run it yourself, all data stored within your control.
* Extensible framework which allows for you to include your own custom enforcements.

# Getting Started

The easiest way to get started is to clone the sample sift repo. It has a
skeleton with the necessary directories for you to get started.

Clone the repo with

   git clone git@github.com:siftproject/sift-repo.git

Update the credentials

   vi sift-repo/accounts/aws/default.json

Run sift

   sift repo -d sift-repo

## Sift Repo Overview

## Policies

The core of a sift repo is it's policies. These are the valiations to perform against
a set of resources. Sift currently supports two types of policies, **verification** and **report**.

### Verifications

Verifications are used to validate specific attributes of the resources in a collection.
Does the instance use the one of the following AMI IDs?  Does a user account have
multi factor authentication enabled?

For example, to validate that all users have MFA enabled, create the following policy file.

    [
      {
        "source" : "aws_iam_user",
        "verifications" : {
          "mfa" : {
            "value" : ["true"]
          }
        }
      }
    ]

By default, sift compares that a value matches, however the following additional
comparisons can be made of value to a desired state.

* include
* exclude
* equals
* within
* set

For example, to ensure that all EC2 instances use one of a specific list of AMIS:

    [
      {
        "source" : "aws_ec2_instance",
        "verifications" : {
          "image" : {
            "comparison" : ["within"],
            "value" : ["ami-12345678", "ami-87654321"]
          }
        }
      }
    ]

### Reports

Reports are ran against an entire collection of resources to make sure it matches a
given set of ids. For example, does the list of users match user1, user2 and user3?
Are there at least 2 snapshots of the database?  Are there less than 20 instances running?

For example, to verify that there are less then 5 instances running, you can use
 the following:

    [
      {
        "source" : "aws_ec2_instance",
        "reports" : {
          "less_than" : ["5"]
        }
      }
    ]

Reports can also validate the entries in a list using **equals**.  For example, to
ensure a list of iam users only contains john and jane:

    [
      {
        "source" : "aws_iam_user",
        "reports" : {
          "equals" : ["john", "jane"]
        }
      }
    ]

Reports can perform the following comparisons.

* equals - The list matches the provided list exactly.
* greater\_then - The list has more than the specified # of entries.
* less\_than - The list has less than the specified # of entries.

## Accounts

Contain credentials which are used to access providers by a given account. Accounts 
will have different credentials depending on the provider.

For example, to access AWS you will need to add an account with a **secret_access_key** and
**access_key_id**. For example:

    {
      "credentials" : {
        "key"    : "abc",
        "secret" : "123"
      }
    }

Accounts can also be referenced by scope, this allows for groups of like accounts to 
be targeted by specific policies. For example, you can put the above account in the
**prod** and **web** scopes.

    {
      "credentials" : {
        "key"    : "abc",
        "secret" : "123"
      },
      "scope" : ["prod", "web"]
    }

A report or verification can then be run against all accounts in a given scope
by referencing it in the verification.

    [
      {
        "source" : "aws_ec2_instance",
        "scope" : "web",
        "verifications" : {
          "image" : {
            "comparison" : ["within"],
            "value" : ["ami-12345678", "ami-87654321"]
          }
        }
      }
    ]

## Sources

Sources provide arguments to providers. For example, to connect to AWS EC2,
you will need to provide the region.

By default, sift will load the default source from the file with the name that
matches the target in **sources**.

For example, to add a source for **aws_ec2_instance**, create the file
**aws_ec2_instance.json** with the following content in sources:

    {
      "default" {
        "region" : ["us-west-1"]
      }
    }

If multiple arguments are provided for the source, sift will run evaluations
against that souce using each of the arguments.

For example, to run a policy against both US West regions, you would specify

    {
      "default" {
        "region" : ["us-west-1", "us-west-2"]
      }
    }

Multiple source arguments can be set by giving them each a unique name which is
referenced in the policy. For example, to set one policy for all us regions, and one
for only the us-east-1 region, you can create the following source:

    {
      "us" {
        "region" : ["us-west-1", "us-west-2", "us-east-1"]
      },
      "east" {
        "region" : ["us-east-1"]
      }
    }

You can specify which arguments to provide in the verificaiton

    [
      {
        "source" : "aws_ec2_instance",
        "arguments" : "us",
        "reports" : {
          "less_than" : ["5"]
        }
      }
    ]

# Advanced

## Repo Overview

The repo directory struct is layed out as follows.

|-accounts
          \aws.json
           github.json
|-filters
         \filter1.json
          filter2.json
|-lists
       \users
             \user1.json
              user2.json
       \images
             \image1.json
              image2.json
|-policies
          \policies_file1.json
           policies_file2.json
|-sources
         \resource1.json
          resource2.json

* The account file name maps to an account for which you are providing credentials.
* Filters files have arbitrary names and contain one or more filters.
* Lists contain directories contain an list of items.
* Policies files have arbitrary names and contain one or more policies.
* Sources map to a given account-provider-collection (eg. aws\_ec2\_\instance).


## Lists

There are times when you have a list of many items to check, or when you have
a resource which may have multiple names across accounts or providers. For example
the user John maybe john123 on github, but jdoe on AWS. Lists allow you to create
a mapping for a single entity, to multiple different names.

## Filters

Filters are used to target a policy at a specific subset of resources in a
collection. A policy can be filtered by the following

* include
* exclude
* attributes

### Attributes

Only include resources with the given attributes. This uses the same syntax as
policy attribute verifications.

#### Layered Filters

Attribute filters can be layered to provide very targeted set of resources.
For example, to only include objects within buckets that have versioning enabled, and
to exclude objects with the ID "123", you can apply the following bucket-object filter.

```json
    "filters" : {
      "bucket" : {
        "attributes" : {
          "versioning_enabled" : {
            "equals" : ["true"]
          }
        }
      },
      "bucket-object" : {
        "exclude" : ["y"]
      }
    },
```

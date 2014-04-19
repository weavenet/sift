# Sift

Sift provides configuration and policy verification and audit fromework for application and organizations running at scale in the cloud.

# Philosophy

Maintaining secuirty and proper configuraltin when running large application
in the cloud is difficult. Enterprises will in a very short
amount of time find that they have dozens of different accounts, across multiple
providers, service many different projects. Those accounts are managed by many
different engineers, across different roles, with varyings levels of experience.

You want the engineers building and running these projects to have control
of their cloud resources, however with the proper guardrails in place to ensure
they are not exposing their products to risk in the cloud.

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

    git clone git@github.com:brettweavnet/sift-repo.git

Build Sift

    make

Update the credentials

    vi sift-repo/accounts/aws/default.json

Run sift against the "loopback" example repo. This repo simulates running a variety
of verifications against a mocked out AWS provider.

    sift repo -d ./examples/repo/passing

## Sift Repo Overview

## Policies

The core of a sift repo is it's policies. These are the valiations to perform against
a resource, or a set of resources. Sift currently supports two types of
policies, **verification** and **report**.

### Verifications

Verifications are used to validate specific attributes of the resources in a collection.
Do instances use the one of the following AMI IDs?  Do user account have multi factor
authentication enabled?

For example, to validate that all users have MFA enabled, create the following policy file.

```json
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
```

By default, sift ensures that the value **equals** the resource state.

The following comparisons of the value and resource state can be made.

* equals
* include
* exclude
* equals
* within
* set

For example, to ensure that all EC2 instances use one of a specific list of AMIS:

```json
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
```

### Reports

Reports are ran against an entire collection of resources to make sure it matches a
list or quantity. For example, does the list of users match user1, user2 and user3?
Are there at least 2 snapshots of the database? Are there less than 20 instances running?

For example, to verify that there are less then 5 instances running, you can use 
the following report:

```json
[
  {
    "source" : "aws_ec2_instance",
    "reports" : {
      "less_than" : ["5"]
    }
  }
]
```

Reports can also validate the entries in a list using **equals**.  For example, to
ensure a list of iam users only contains john and jane:

```json
[
  {
    "source" : "aws_iam_user",
    "reports" : {
      "equals" : ["john", "jane"]
    }
  }
]
```

Reports can perform the following comparisons.

* equals - The list matches the provided list exactly.
* greater\_then - The list has more than the specified # of entries.
* less\_than - The list has less than the specified # of entries.

## Accounts

Accounts contain credentials which are used to access providers by a given account. Accounts
will have different credentials depending on the provider.

Sift repos are expected to be checked into your SCM. The best way to set credentials
is via environment variables which can be set using a secure process.

For example, to access AWS you will need to add an account with a **secret_access_key** 
and **access_key_id**.

```json
{
  "name" : "my_aws_account",
  "credentials" : {
    "access_key_id"     : "$AWS_ACCESS_KEY_ID",
    "secret_access_key" : "$AWS_SECRET_ACCESS_KEY"
  }
}
```

Accounts can also be referenced by scope, this allows for groups of like accounts to 
be targeted by specific policies. For example, you can put the above account in the
**prod** and **web** scopes.

```json
{
  "name" : "my_aws_account",
  "credentials" : {
    "key"    : "abc",
    "secret" : "123"
  },
  "scope" : ["prod", "web"]
}
```

A report or verification can then be run against all accounts in a given scope
by referencing it in the verification.

```json
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
```

As the number of teams and projects grows, accounts can be segrated into different
scopes (Prod / Dev, Web / App, Finance / Marketing, etc) to ensure policies are targeted
correctly.

## Sources

Sources provide arguments to providers. For example, to connect to AWS EC2,
you will need to provide the region.

By default, sift will load the default source from the file with the name that
matches the target in **sources**.

For example, to add a source for **aws_ec2_instance**, create the file
**aws_ec2_instance.json** with the following content in sources:

```json
{
  "default" {
    "region" : ["us-west-1"]
  }
}
```

If multiple arguments are provided for the source, sift will run evaluations
against that souce using each of the arguments.

For example, to run a policy against both US West regions, you would specify

```json
{
  "default" {
    "region" : ["us-west-1", "us-west-2"]
  }
}
```

Multiple source arguments can be set by giving them each a unique name which is
referenced in the policy. For example, to set one policy for all us regions, and one
for only the us-east-1 region, you can create the following source:

```json
{
  "us" {
    "region" : ["us-west-1", "us-west-2", "us-east-1"]
  },
  "east" {
    "region" : ["us-east-1"]
  }
}
```

You can specify which to use in a given policy via the **arguments** key.

```json
[
  {
    "source" : "aws_ec2_instance",
    "arguments" : "us",
    "reports" : {
      "less_than" : ["5"]
    }
  }
]
```

# Advanced

## Repo Overview

The repo directory is layed out as follows.

```
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
         \provider_collection_resource1.json
          provider_collection_resource2.json
```

* The account file name maps to an account for which you are providing credentials.
* Filters files have arbitrary names and contain one or more filters.
* Lists contains directories which are mapped to lists (list users has user1
and user2 in above example).
* Policies files have arbitrary names and contain one or more policies.
* Sources map to a given account-provider-collection (eg. aws\_ec2\_\instance).

## Lists & Functions

There are times when you have a list of many resources to check across multiple
providers, for example your organizations list of valid users accounts.

Additionally, these resources may have multiple names across accounts or providers.
For example the user John Doe maybe johndoe on one aws account, but jdoe on another AWS.
Lists allow you to create a mapping for a single entity, to multiple different names.

Lists and functions provide a way to accomplish both.

Each directory under the lists directory is the name of a list. The individual json
files within the directory are entries in the list.

For example, to create a users list with johndoe and janedoe, create the directory **users** and the files **johndoe.json** and **janedoe.json** under the lists directory.

The following function can be used to insert that user list into a policy.

```json
{ "Fn::List" : ["users"] }
```

Users will have different names across different services, to add an alias to the user,
update the json file to include the primary id, as well as any aliases.  For example
to add an alias for the finance account for **johndoe**:

```json
{
  "id": "johndoe",
  "finance": "jdoe1"
}
```

The following code will then over-ride the primary id with finance, when available.

```json
{ "Fn::ListSub" : ["users", "finance"] }
```

Thie will result in the following users.

```json
["jdoe1", "janedoe"]
```

The **ListOnly** function can be used to substitute only those users who are in the finance group.

```json
{ "Fn::ListOnly" : ["users", "finance"] }
```

Will result in

```json
["jdoe1"]
```

## Filters

Filters are used to target a policy at a specific subset of resources in a
collection. A policy can be filtered by the following

* include - including specified resource IDs
* exclude - excluding specified resource IDs
* attributes - including only resources whose attributes match the given criteria

The filters directory contains files, with arbitrary names, that contain filters
which can be referenced within policies.

For example, to create a filter to only include web instances, create a
json file with the following content under the filters directory of the repo:

```json
{
  "web" : {
    "include" : ["i-1234abcd", "i-9876abcd"]
  }
}
```

To create another filter exclude another set of instances:

```json
{
  "exclude_db" : {
    "exclude" : ["i-abcd1234"]
  }
}
```

This can then be referenced in a policy via:

```json
{
  "filters": {
    "instance": "exclude_db"
  },
  "source" : "aws_ec2_instance",
  "reports" : {
    "less_than" : ["3"]
  }
}
```

### Attributes

Only include resources with the given attributes. This uses the same syntax as
policy attribute verifications.

#### Layered Filters

Attribute filters can be layered to provide very targeted set of resources.
For example, to only include objects within buckets that have versioning enabled, and
to exclude objects with the ID "123", you can apply the following bucket-object filter.

```json
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
```

# Known Issues

* Attribute filtering current has bugs.

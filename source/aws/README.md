# AWS

## Ec2

### instance

#### Attributes

* id
* image_id

#### Example

```json
{
  "source" : "aws_ec2_instance",
  "verifications" : {
    "image_id" : {
      "comparison" : ["within"],
      "value" : ["ami-12345678", "ami-87654321"]
    }
  }
}
```

### securitygroup

#### Attributes

* id
* name
* vpc_id

#### Example

```json
{
  "source" : "aws_ec2_securitygroup",
  "verifications" : {
    "vpc_id" : {
      "value" : ["vpc-12345678"]
    }
  }
}
```

### securitygroup-ippermission

#### Attributes

* protocol
* to_port
* from_port
* source_ips

#### Example

```json
{
  "source" : "aws_ec2_securitygroup-ippermission",
  "verifications" : {
    "source_ips" : {
      "comparison" : ["exclude"],
      "value" : ["0.0.0.0/0"]
    },
    "protocol" : {
      "comparison" : ["within"],
      "value" : ["tcp", "udp"]
    }
  }
}
```

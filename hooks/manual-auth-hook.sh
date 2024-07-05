#!/bin/bash

echo "Please deploy a DNS TXT record under the name:"
echo "_acme-challenge.$CERTBOT_DOMAIN."
echo "with the following value:"
echo "$CERTBOT_VALIDATION"

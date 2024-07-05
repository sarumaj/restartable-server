#!/bin/bash

echo "Please deploy a DNS TXT record under the name:"
echo "_acme-challenge.$CERTBOT_DOMAIN."
echo "with the following value:"
echo "$CERTBOT_VALIDATION"

# Output the DNS record to the GitHub Actions log for manual addition
echo "::set-output name=dns_name::_acme-challenge.$CERTBOT_DOMAIN"
echo "::set-output name=dns_value::$CERTBOT_VALIDATION"

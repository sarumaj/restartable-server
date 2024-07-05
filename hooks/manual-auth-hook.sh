#!/bin/bash

# Set the environment variables for Namecheap API
API_USER=${NAMECHEAP_API_USER}
API_KEY=${NAMECHEAP_API_KEY}
USER_NAME=${NAMECHEAP_API_USER}
CLIENT_IP=${NAMECHEAP_API_IP}
DOMAIN=${DOMAIN}

# Extract the SLD (second-level domain) and TLD (top-level domain)
SLD=$(echo $DOMAIN | cut -d'.' -f1)
TLD=$(echo $DOMAIN | cut -d'.' -f2-)

# Display the DNS TXT record that needs to be added
echo "Please deploy a DNS TXT record under the name:"
echo "_acme-challenge.$CERTBOT_DOMAIN."
echo "with the following value:"
echo "$CERTBOT_VALIDATION"

# Output the DNS record to the GitHub Actions log for manual addition
echo "::set-output name=dns_name::_acme-challenge.$CERTBOT_DOMAIN"
echo "::set-output name=dns_value::$CERTBOT_VALIDATION"

# Get public IP
IP=$(curl -s ifconfig.me)
echo "Public IP: $IP"
echo "::set-output name=public_ip::$IP"

# Define the authentication string
AUTH="ApiUser=$API_USER&ApiKey=$API_KEY&UserName=$USER_NAME&ClientIp=$CLIENT_IP"

# Fetch existing DNS records
HOSTS=$(curl -s "https://api.namecheap.com/xml.response?${AUTH}&Command=namecheap.domains.dns.getHosts&SLD=${SLD}&TLD=${TLD}")

# Prepare the XML for the existing hosts
HOSTS_XML=""
IFS=$'\n' # Make sure the loop handles lines properly

# Parse existing hosts and add them to the XML string
echo "$HOSTS" | grep -oPm1 "(?<=<Host ).*?(?=\</Host>)" | while read -r HOST; do
    NAME=$(echo "$HOST" | grep -oPm1 "(?<=Name=\").*?(?=\")")
    TYPE=$(echo "$HOST" | grep -oPm1 "(?<=Type=\").*?(?=\")")
    ADDRESS=$(echo "$HOST" | grep -oPm1 "(?<=Address=\").*?(?=\")")
    MX_PREF=$(echo "$HOST" | grep -oPm1 "(?<=MXPref=\").*?(?=\")")
    TTL=$(echo "$HOST" | grep -oPm1 "(?<=TTL=\").*?(?=\")")

    HOSTS_XML="${HOSTS_XML}<Host Name=\"$NAME\" Type=\"$TYPE\" Address=\"$ADDRESS\" MXPref=\"$MX_PREF\" TTL=\"$TTL\"/>"
done

# Add the new TXT record for Let's Encrypt validation
HOSTS_XML="${HOSTS_XML}<Host Name=\"_acme-challenge\" Type=\"TXT\" Address=\"$CERTBOT_VALIDATION\" TTL=\"60\"/>"

# Update the DNS records with the new XML
UPDATE_RESULT=$(curl -s "https://api.namecheap.com/xml.response?${AUTH}&Command=namecheap.domains.dns.setHosts&SLD=${SLD}&TLD=${TLD}" \
    --data-urlencode "HostNames=${HOSTS_XML}")

# Output the result
echo "$UPDATE_RESULT"

# Check if the update was successful
if echo "$UPDATE_RESULT" | grep -q "<ErrCount>0</ErrCount>"; then
    echo "DNS records updated successfully."
else
    echo "Failed to update DNS records."
    exit 1
fi

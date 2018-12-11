#!/bin/zsh


address=`http -b POST buaa01:8090/api/ipaddress name=php-apache | grep ipaddress | jq .ipaddress | sed 's/\"//g'`
echo $address

while (( 1 )) {
    curl $address
}

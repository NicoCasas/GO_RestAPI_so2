#!/bin/bash

# Ejecuta todos los tests

echo "test_login_fail : $(./test_login_fail.sh)"
echo "test_login_pass : $(./test_login_pass.sh)"
echo "test_get_users  : $(./test_get_users.sh)"
echo "test_create_user: $(./test_create_user.sh)"
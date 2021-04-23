eval `ssh-agent`

ssh-add /root/.ssh/id_rsa

ansible-playbook -i ansible_hosts provision.yml --private-key=/data/certs/balconygames.pem --extra-vars "@group_vars/prod.secure.yml" --extra-vars "@group_vars/prod.yml" --vault-password-file /data/certs/balconygames.pw

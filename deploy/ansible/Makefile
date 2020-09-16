encrypt:
	pipenv run ansible-vault encrypt --vault-password-file ~/.balconygames.pw group_vars/prod.secure.yml
.PHONY: encrypt

decrypt:
	pipenv run ansible-vault decrypt --vault-password-file ~/.balconygames.pw group_vars/prod.secure.yml
.PHONY: decrypt

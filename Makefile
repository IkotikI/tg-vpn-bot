

# ---- Admin panel commands ---- 

# Start Golang Air for admin_panel directory.
admin: air-admin
# Compile CSS, JS and Golang Templ. It used in Air config. 
admin-prepare: admin-tailwind admin-typescript templ-generate

admin-old:
	templ generate
	cd web/admin_panel/ && $(MAKE) tailwind
	air -c cmd/admin/.air.toml
	# cd cmd/admin/ && air -c .air.toml
	# go run cmd/admin/admin.go

admin-tailwind:
	@cd web/admin_panel/ && $(MAKE) tailwind

admin-tailwind-watch:
	@cd web/admin_panel/ && $(MAKE) tailwind-watch

admin-typescript:
	@cd web/admin_panel/ && $(MAKE) typescript

templ-generate:
	templ generate

air-admin:
	air -c cmd/admin/.air.toml




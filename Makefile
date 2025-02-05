admin-old:
	templ generate
	cd web/admin_panel/ && $(MAKE) tailwind
	air -c cmd/admin/.air.toml
	# cd cmd/admin/ && air -c .air.toml
	# go run cmd/admin/admin.go

tailwind:
	cd web/admin_panel/ && $(MAKE) tailwind

tailwind-watch:
	cd web/admin_panel/ && $(MAKE) tailwind-watch

templ-generate:
	templ generate

air-admin:
	air -c cmd/admin/.air.toml


admin: air-admin

.PHONY: mailhog
mailhog:
	docker run -p 1025:1025 -p 8025:8025 mailhog/mailhog
include META

.PHONY: swag
swag:
	@echo "[GENERATE DOCS] Generating API documents"
	@echo " - Initializing swag"
	@swag init --parseDependency --parseInternal --generatedTime

.PHONY: push
push:
	@echo "[BUILD] Committing and pushing to remote repository"
	@echo " - Committing"
	@git add META
	@git commit -am "v$(VERSION)"
	@echo " - Tagging"
	@git tag v${VERSION}
	@echo " - Pushing"
	@git push --tags origin ${BRANCH}
BRANCH = $(shell git symbolic-ref --short HEAD)

OS = ${shell uname -s}
ifeq (${OS},Linux)
	SED_OPT := -i
else ifeq (${OS},Darwin)
	SED_OPT := -i ''
endif

expect-%:
	@test "$(BRANCH)" = "$*" || \
		(echo "ERROR: not in $* branch" && exit 1)

merge-gitlab: expect-master
	test "`git stash`" = 'No local changes to save' && STASH=0 || STASH=1; \
		git checkout -B gitlab && \
		git merge $(BRANCH) && \
		sed $(SED_OPT) 's|github.com/snobb/susanin|gitlab.com/snobb/susanin|' ./pkg/*/*.go ./examples/*.go go.mod README.md && \
		git commit -m 'merging to gitlab' -a && \
		git push gitlab gitlab -f && \
		git checkout master && \
		[ $$STASH -eq 1 ] && git stash pop

merge-github: expect-master
	test "`git stash`" = 'No local changes to save' && STASH=0 || STASH=1; \
		git checkout -B github && \
		git merge $(BRANCH) && \
		sed $(SED_OPT) 's|github.com/snobb/susanin|github.com/snobb/susanin|' ./pkg/*/*.go ./examples/*.go go.mod README.md && \
		git commit -m 'merging to github' -a && \
		git push github github -f && \
		git checkout master && \
		[ $$STASH -eq 1 ] && git stash pop

merge-all: merge-gitlab merge-github

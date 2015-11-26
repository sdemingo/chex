
files = chex-users.js chex-questions.js chex-answers.js chex-panellist.js chex-tests.js chex-utils.js chex-main.js 
target = chex.js
target_min = chex.min.js

$(target): $(files)
	echo "Building $(target)"
	echo > $(target)
	for file in $(files) ; do \
		cat $$file >> $(target) ; \
	done
	#uglifyjs $(target) > $(target_min)

all: $(target)

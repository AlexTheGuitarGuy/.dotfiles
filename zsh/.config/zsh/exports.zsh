# Setting for the new UTF-8 terminal support in Lion
export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

export GEM_HOME="$(ruby -e 'puts Gem.user_dir')"
export PATH="$PATH:$GEM_HOME/bin"
export PATH="$PATH:/home/alex/.cargo/bin"

PATH="/home/alex/perl5/bin${PATH:+:${PATH}}"
export PATH
PERL5LIB="/home/alex/perl5/lib/perl5${PERL5LIB:+:${PERL5LIB}}"
export PERL5LIB
PERL_LOCAL_LIB_ROOT="/home/alex/perl5${PERL_LOCAL_LIB_ROOT:+:${PERL_LOCAL_LIB_ROOT}}"
export PERL_LOCAL_LIB_ROOT
PERL_MB_OPT="--install_base \"/home/alex/perl5\""
export PERL_MB_OPT
PERL_MM_OPT="INSTALL_BASE=/home/alex/perl5"
export PERL_MM_OPT
export PATH="/usr/pgsql-14/bin:$PATH"

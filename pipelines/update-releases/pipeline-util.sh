#!/bin/bash -eu
## ======================================================================
## ./pipeline-util.sh --opt_test
## ./pipeline-util.sh --opt_gen
## ======================================================================

## Prerequisites
type texplate > /dev/null
type bosh > /dev/null

opt_test=""
opt_gen=""
opt_destroy=""
opt_sort=""
opt_psdiff=""

while [[ "$#" > 0 ]]; do
    case $1 in
        --opt_test)
            opt_test="true"
            shift
            ;;
        --opt_gen)
            opt_gen="true"
            shift
            ;;
        --opt_destroy)
            opt_destroy="true"
            shift
            ;;
        --opt_sort)
            opt_sort="true"
            shift
            ;;
        --opt_psdiff)
            opt_psdiff="true"
            shift
            ;;
        *) break
           ;;
  esac;
done

if [ -n "${opt_test}" ]; then
	## Prerequisites
	# old_job=$(mktemp)
	# new_job=$(mktemp)
	old_job=update-orig.yml
	new_job=update-new.yml
	releases=${releases:="
                          backup-and-restore-sdk
                          bosh-dns
                          bpm
                          broker-registrar
                          cf-app-sd
                          credhub
                          haproxy
                          nfs-volume
                          postgres
                          silk
                          syslog
                          
                          binary-buildpack
                          capi
                          cf-mysql
                          cf-networking
                          cf-smoke-tests
                          cf-syslog-drain
                          cflinuxfs2
                          consul
                          diego
                          dotnet-core-buildpack
                          garden-runc
                          go-buildpack
                          java-buildpack
                          loggregator
                          nats
                          nodejs-buildpack
                          php-buildpack
                          python-buildpack
                          routing
                          ruby-buildpack
                          staticfile-buildpack
                          statsd-injector
                          uaa
                  
                          windows-stemcell
                          windows2016-stemcell
                         "}

	for release in $releases; do
	    set +e
            compare_path=/jobs/name=update-$release
     	    bosh int ../update-releases.yml --path ${compare_path} > $old_job

            compare_path=/jobs/name=update-$release
	        bosh int <(texplate execute template.yml --input-file input.yml) --path ${compare_path} > $new_job

            if [ -n "${opt_sort}" ]; then
                sort ${old_job} > ${old_job}.sorted
                sort ${new_job} > ${new_job}.sorted
	            diff_output=$(diff ${old_job}.sorted ${new_job}.sorted)
            else
	            diff_output=$(diff ${old_job} ${new_job})
            fi
	        exit_code=$?

	    set -e

	    if [ $exit_code != 0 ]; then
	        echo "update-$release jobs have diverged: $diff_output"
	        exit 2
	    else
	        echo "update-$release jobs are in sync"
	        rm -f $old_job $new_job
	    fi
	done
elif [ -n "${opt_destroy}" ]; then
    echo "Delete pipeline: update-releases-gen"
    fly --target glenda-ci \
        destroy-pipeline \
        --pipeline=update-releases-gen \
        --non-interactive
    fly --target glenda-ci \
        destroy-pipeline \
        --pipeline=update-releases-original \
        --non-interactive
elif [ -n "${opt_gen}" ]; then
    echo "Generate pipeline: update-releases-gen.yml"
	texplate execute template.yml --input-file input.yml > update-releases-gen.yml
    echo "Set pipeline: update-releases-gen.yml"
    fly --target glenda-ci \
        set-pipeline \
        --config=update-releases-gen.yml \
        --pipeline=update-releases-gen \
        --load-vars-from=$HOME/workspace/runtime-ci-private/pipeline_vars/update-releases.yml \
        --non-interactive > /dev/null
    echo "Set pipeline: ../update-releases.yml"
    fly --target glenda-ci \
        set-pipeline \
        --config=../update-releases.yml \
        --pipeline=update-releases-original \
        --load-vars-from=$HOME/workspace/runtime-ci-private/pipeline_vars/update-releases.yml \
        --non-interactive > /dev/null
elif [ -n "${opt_psdiff}" ]; then
    fly --target glenda-ci \
        get-pipeline \
        --pipeline=update-releases-original \
        > update-releases-get-pipeline.yml
    
    sort update-releases-get-pipeline.yml > update-releases-get-pipeline.yml.sorted
    
    fly --target glenda-ci \
        get-pipeline \
        --pipeline=update-releases-gen \
        > update-releases-gen-get-pipeline.yml
    
    sort update-releases-gen-get-pipeline.yml > update-releases-gen-get-pipeline.yml.sorted

    diff update-releases-get-pipeline.yml update-releases-gen-get-pipeline.yml

fi


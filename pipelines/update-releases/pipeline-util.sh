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
opt_pipeline_diff=""

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
        --opt_pipeline_diff)
            opt_pipeline_diff="true"
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
	releases=${releases:="haproxy
	          broker-registrar
	          backup-and-restore-sdk
	          nfs-volume
	          binary-buildpack
	          capi
	          cf-app-sd
	          cf-networking
	          silk
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
	          cf-mysql
	          nats
	          nodejs-buildpack
	          php-buildpack
	          postgres
	          python-buildpack
	          routing
	          ruby-buildpack
	          staticfile-buildpack
	          statsd-injector
	          uaa
	          windows-stemcell
	          windows2016-stemcell
	         "}

    releases=credhub
	for release in $releases; do
	    set +e

	        if [ $release = "backup-and-restore-sdk" ]; then
	            bosh int ../update-releases.yml --path /jobs/name=update-bbr > $old_job
	            sed -i .bak 's/bbr/backup-and-restore-sdk/g' $old_job
	            rm -f $old_job.bak
	        else
                compare_path=/jobs/name=update-$release
	            bosh int ../update-releases.yml --path ${compare_path} > $old_job
	        fi

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
elif [ -n "${opt_gen}" ]; then
    echo "Generate pipeline: update-releases-gen.yml"
	texplate execute template.yml --input-file input.yml > update-releases-gen.yml
    echo "Set pipeline: update-releases-gen.yml"
    fly --target glenda-ci \
        set-pipeline \
        --config=update-releases-gen.yml \
        --pipeline=update-releases-gen \
        --load-vars-from=$HOME/workspace/runtime-ci-private/pipeline_vars/update-releases.yml \
        --non-interactive
elif [ -n "${opt_pipeline_diff}" ]; then
    fly380 --target relint-ci \
        get-pipeline \
        --pipeline=update-releases \
        > update-releases-get-pipeline.yml
    fly --target glenda-ci \
        get-pipeline \
        --pipeline=update-releases-gen \
        > update-releases-gen-get-pipeline.yml

    diff update-releases-get-pipeline.yml update-releases-gen-get-pipeline.yml

fi


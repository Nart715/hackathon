#!/bin/sh

ports[0]=7001
ports[1]=7002
ports[2]=7003
ports[3]=7004
ports[4]=7005
ports[5]=7006

root_dir=$(pwd)
echo "CURRENT ROOT DIR: $root_dir"
redis_config_file_name="redis.conf"
redis_ver="7.2.2"
compress_type="tar.gz"
redis_dir_name="redis-$redis_ver"


# copy to staging: scp -P 2115 -r *  sysadmin@221.132.18.243:taitran/redis
# copy to container dir: docker cp . 6448491334cd:/code/redis
# start this

create_dir(){
	echo ">>>>> CREATE DIR <<<<<"
	 for i in "${ports[@]}"; do
	 	cd $root_dir
		echo "CREATE DIR $i\n"
	 	mkdir $i
	 done
}

download_redis(){
	wget https://download.redis.io/redis-stable.tar.gz
}

create_dir_folder(){
	cd /code
	mkdir 7001 7002 7003 7004 7005 7006
	cp -r $redis_dir_name 7001/
	cp -r $redis_dir_name 7002/
	cp -r $redis_dir_name 7003/
	cp -r $redis_dir_name 7004/
	cp -r $redis_dir_name 7005/
	cp -r $redis_dir_name 7006/
}

install_redis_local_extract(){
	cd $root_dir/7001
	cd $redis_dir_name
	make

	cd $root_dir/7002
	cd $redis_dir_name
	make

	cd $root_dir/7003
	cd $redis_dir_name
	make

	cd $root_dir/7004
	cd $redis_dir_name
	make

	cd $root_dir/7005
	cd $redis_dir_name
	make

	cd $root_dir/7006
	cd $redis_dir_name
	make

}

install_redis_cluster_link(){
	cd $root_dir/7001
	wget https://download.redis.io/redis-stable.tar.gz
	tar xzf redis-stable.tar.gz
	cd redis-stable
	make

	cd $root_dir/7002
	wget https://download.redis.io/redis-stable.tar.gz
	tar xzf redis-stable.tar.gz
	cd redis-stable
	make

	cd $root_dir/7003
	wget https://download.redis.io/redis-stable.tar.gz
	tar xzf redis-stable.tar.gz
	cd redis-stable
	make

	cd $root_dir/7004
	wget https://download.redis.io/redis-stable.tar.gz
	tar xzf redis-stable.tar.gz
	cd redis-stable
	make

	cd $root_dir/7005
	wget https://download.redis.io/redis-stable.tar.gz
	tar xzf redis-stable.tar.gz
	cd redis-stable
	make

	cd $root_dir/7006
	wget https://download.redis.io/redis-stable.tar.gz
	tar xzf redis-stable.tar.gz
	cd redis-stable
	make
}

install_redis_cluster_package_local(){
	cd $redis_dir_name
	make
}

install_redis_packages(){
	echo " >>>> INSTALL REDIS PACKAGES <<<<"
	tar xzf $redis_dir_name.$compress_type
	for i in "${ports[@]}"; do
		echo "$redis_dir_name.$compress_type"
    echo "COPY $redis_dir_name TO $root_dir/$i"
    cp -r $redis_dir_name $root_dir/$i
	 	cd $root_dir/$i
    pwd
		echo "INSTALL REDIS $i\n"
	 	install_redis_cluster_package_local
    cd $root_dir
	done
	echo "\n"
}

copy_config_redis_cluster(){
	echo ">>>>> CONFIG REDIS CLUSTER <<<<<"
	for i in "${ports[@]}"; do
		echo "COPY CONFIG REDIS $i\n"
	 	cp -r "$root_dir/redis_$i.conf" $root_dir/$i/$redis_dir_name/$redis_config_file_name
	done
	echo "\n"
}

start_redis_host(){
	echo ">>>>> START REDIS HOST <<<<<"
	for i in "${ports[@]}"; do
		cd $root_dir/$i/$redis_dir_name
		echo "START HOST REDIS $i\n"
	 	./src/redis-server redis.conf &
	done
	echo "\n"
}

sync_redis_cluster(){
	echo ">>>>> SYNC REDIS CLUSTER <<<<<"
	# not index meaning to get the first index = 0
	echo "$ports"
	cd $root_dir/$ports/$redis_dir_name
	pwd
	clause="src/redis-cli --cluster create"
	for i in "${ports[@]}"; do
		clause="${clause} 127.0.0.1:$i "
	done
	clause="${clause} --cluster-replicas 1"
	echo $clause
  ./$clause
}

export_redis_cli(){
	alias redis_cli='cd /code/7001/redis-5.0.5/src/ && ./redis-cli -h 127.0.0.1 -p 7001 -c'
}

create_dir
install_redis_packages
copy_config_redis_cluster
start_redis_host
sync_redis_cluster

# ref(unicorn設定): https://qiita.com/syou007/items/555062cc96dd0b08a610
# ref(unicorn/h2o設定): https://orihubon.com/2015/08/19/h2o/

# project-root/config/unicorn.rb
@dir = File.expand_path('../../', __FILE__)

worker_processes 1 # CPUコア数によって変えられるそう

working_directory @dir

timeout 30

stderr_path File.expand_path('../../log/unicorn_stderr.log', __FILE__)
stdout_path File.expand_path('../../log/unicorn_stdout.log', __FILE__)

#listen 8777, :tcp_nopush => true
#listen File.expand_path('../../tmp/sockets/unicorn.sock', __FILE__)
listen File.expand_path('/tmp/sockets/unicorn.sock', __FILE__)

pid File.expand_path('../../tmp/pids/unicorn.pid', __FILE__)

preload_app true

# 以下ファイルにパーミッションを与えること
# unicorn_stderr.log
# unicorn_stdout.log
# unicorn.pid
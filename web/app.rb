require 'sinatra'
require 'sinatra/base'
require 'sinatra/reloader'


get '/' do
  'sinatra!!!!!'
end

get '/tmp_erb' do
  erb :tmp_erb
end

get '/list' do
  @foo = 'get_list'
  erb :'/memo/list'
end

# client error
error 400..499 do
  #status = response.status
  @e = "クライアントエラー: #{response.status}"
  @msg = "正しい操作をしてください。"
  erb :error
end

# server error
error 500..599 do
  @e = "サーバーエラー: #{response.status}"
  @msg = "管理者に連絡してください。"
  erb :error
end

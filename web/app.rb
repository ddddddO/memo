require 'sinatra'
require 'sinatra/base'
require 'sinatra/reloader'

require './models/model'

# ref:sinatra を要確認 

configure do
  set(:model) { Model.new } # DB接続処理
end

get '/' do
  erb :auth
end

enable :sessions

post '/login' do
  user_id = settings.model.login(params[:name], params[:passwd])
  p user_id
  session[:user_id] = user_id
  redirect to('/list')
end

get '/list' do
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

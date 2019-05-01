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
  session[:user_id] = user_id
  redirect to('/list')
end

get '/list' do
  @memos = settings.model.list(session[:user_id])
  erb :'/memo/list'
end

# ディレクトリトラバーサルへは、sessionのuser_idを条件に含めSQL実行で回避
# TODO: 不正なパスを入力された場合のハンドリング
get '/detail/:memo_id' do
  @memo = settings.model.detail(params[:memo_id], session[:user_id])
  erb :'/memo/detail'
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

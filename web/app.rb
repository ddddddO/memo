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
# sessionとcookieの関係については以下
# https://riocampos-tech.hatenablog.com/entry/20140616/private_study_about_rack_session_or_cookie
# github.com/ddddddO/work/ruby/dec_cookie_in_session.rb
#
# set :sessions, secret: 'xxx'

post '/login' do
  user_id = settings.model.login(params[:name], params[:passwd])
  session[:user_id] = user_id

  redirect to('/list')
end

get '/list' do
  if session[:user_id].nil?
    redirect to('/')
  end

  @memos = settings.model.list(session[:user_id])
  erb :'/memo/list'
end

# TODO: 不正なパスを入力された場合のハンドリング
get '/detail/:memo_id' do
  if session[:user_id].nil?
    redirect to('/')
  end

  @memo = settings.model.detail(params[:memo_id], session[:user_id])
  erb :'/memo/detail'
end

# put method ref: https://qiita.com/suin/items/d17bdfc8dba086d36115
# formのmethodをdelete/putへ対応: http://portaltan.hatenablog.com/entry/2015/07/22/122031
# 更新画面への遷移用
post '/update_view' do
  if session[:user_id].nil?
    redirect to('/')
  end

  if params.nil?
    # メモ新規作成用
    params[:memo_id] = ''
    params[:subject] = ''
    params[:content] = ''
  end

  erb :'memo/update'
end

put '/update' do
  if session[:user_id].nil?
    redirect to('/')
  end

  params[:user_id] = session[:user_id]
  memo_id = settings.model.update(params)

  # メモ詳細へ戻る
  redirect to("/detail/#{memo_id}")
end


# client error
error 400..499 do
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

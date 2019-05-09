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
  if session[:user_id].nil?
    redirect to('/')
  end

  @memos = settings.model.list(session[:user_id])
  erb :'/memo/list'
end

# ディレクトリトラバーサルへは、sessionのuser_idを条件に含めSQL実行で回避
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

  # メモ新規・編集処理はupsertで対応。一旦、DBのid連番対応するまでupdateのみの実装
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

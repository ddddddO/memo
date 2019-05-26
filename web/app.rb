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

  @memos, @tags = settings.model.detail(params[:memo_id], session[:user_id])

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
    params[:tags_length_of_memo] = 0
    params[:tag_ids_of_memo] = []
    params[:tag_names_of_memo] = []
    params[:content] = ''
  end

  # タグ用
  # TODO: いい方法あとで
  if params[:tags_length_of_memo] != 0
    tag_ids_of_memo = []
    tag_names_of_memo = []
    params[:tags_length_of_memo] = params[:tags_length_of_memo].to_i
    params[:tags_length_of_memo].times do |i|
      tag_ids_of_memo.push(params[:"tag_id_#{i}"])
      tag_names_of_memo.push(params[:"tag_name_#{i}"])
    end

    params[:tag_ids_of_memo] = tag_ids_of_memo
    params[:tag_names_of_memo] = tag_names_of_memo
  end

  # ユーザーがもつタグをすべて取得
  params[:all_tags_of_user] = settings.model.fetch_all_tags_of_user(session[:user_id])

  erb :'memo/update'
end

put '/update' do
  if session[:user_id].nil?
    redirect to('/')
  end

  params[:user_id] = session[:user_id]
    
  # タグ(delete)用
  delete_tag_ids = []
  if params.key?('delete_tags_length')
    params[:delete_tags_length].to_i.times do |i|
      if !params.key?("delete_tag_id_#{i}")
        next
      end
      delete_tag_ids.push(params[:"delete_tag_id_#{i}"])
    end
    params[:delete_tag_ids] = delete_tag_ids
  end

  # タグ(update)用
  update_tag_ids = []
  if params.key?('update_tags_length')
    params[:update_tags_length].to_i.times do |i|
      if !params.key?("update_tag_id_#{i}")
        next
      end
      update_tag_ids.push(params[:"update_tag_id_#{i}"])
    end
    params[:update_tag_ids] = update_tag_ids
  end

  p 'paramsss'
  p params

  #memo_id = settings.model.update(params)

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

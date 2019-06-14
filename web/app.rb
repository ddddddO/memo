require 'sinatra'
require 'sinatra/base'
require 'sinatra/reloader'

require './models/model'
require './lib/password'

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
  secure_passwd = Password.gen_secure_password(params[:passwd], params[:name], 100000)
  user_id = settings.model.login(params[:name], secure_passwd)
  session[:user_id] = user_id

  redirect to('/list?page=1')
end

get '/list' do
  if session[:user_id].nil?
    redirect to('/')
  end

  # タグidがクエリパラメータとして送られた場合
  tag_id = ''
  if params.key?('tag_id')
    tag_id = params[:tag_id]
  end
  @memos = settings.model.list(session[:user_id], tag_id) # TODO: rubyでオプション的なのあったからそれ使うようにすれば？
  
  if !params.key?('page')
    params[:page] = '1'
  end

  erb :'/memo/list'
end

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
  # TODO: key名もそれっぽいのに変える
  params[:all_tags_of_user] = settings.model.fetch_all_tags_of_user_excluded_binded_tags(session[:user_id], params['memo_id'])

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

  memo_id = settings.model.update(params)

  # メモ詳細へ戻る
  redirect to("/detail/#{memo_id}")
end

get '/tag_view' do
  if session[:user_id].nil?
    redirect to('/')
  end
  
  @tags = settings.model.tags(session[:user_id])

  erb :'tag/tags'
end

get '/tags_select_for_update' do
  if session[:user_id].nil?
    redirect to('/')
  end

  @tags = settings.model.tags(session[:user_id])
  @tags = @tags[1, @tags.length] # 最初の要素ALLは含めない

  erb :'tag/tags_select_for_update'
end

get '/tag_update_view' do
  if session[:user_id].nil?
    redirect to('/')
  end

  @tag = settings.model.tag(params[:tag_id])

  erb :'tag/update'
end

post '/tag_update' do
  if session[:user_id].nil?
    redirect to('/')
  end
  
  settings.model.update_tag(params)
  redirect to('/list')
end

delete '/tag_delete' do
  if session[:user_id].nil?
    redirect to('/')
  end

  settings.model.delete_tag(params[:tag_id])
  redirect to('/list')
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

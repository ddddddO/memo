require 'pg'

class Model
  # pg: http://www.ownway.info/Ruby/pg/about
  def initialize()
    puts 'connect postgres'
    
    @conn = PG::connect(
        host:     'localhost',
        user:     'postgres',
        password: 'postgres',
        dbname:   'tag-mng'
    )
  end

  def login(name, passwd)
    q = 'SELECT id FROM users WHERE name=$1 AND passwd=$2'
    rslt = @conn.exec(q, [name, passwd])
    
    # TODO: 要リファクタ&認証エラーハンドル(401)
    if !rslt.nil?
      user_id = ''
      rslt.each do |row|
        user_id = row['id']
      end

      if !user_id.empty?
        return user_id
      end
    else
      return 'failed to login'
    end
    return 'failed to login'
  end

  def list(user_id)
    q = 'SELECT id, subject FROM memos WHERE users_id=$1 ORDER BY id'
    rslt = @conn.exec(q, [user_id])

    rows = []
    rslt.each do |row|
      rows.push(row)
    end
    rows
  end

  # タグは別個のSQLでidとnameを取得してrowに詰めて渡した方がいいかも
  def detail(memo_id, user_id)  
=begin
    q = <<~EOS
      SELECT DISTINCT 
        m.id AS id,
        m.subject AS subject,
        m.content AS content,
        ARRAY(
          SELECT 
            t.name
          FROM tags t JOIN memo_tag mtg 
          ON t.id = mtg.tags_id
          WHERE mtg.memos_id = $1
        ) AS tag_names 
      FROM memos m JOIN memo_tag mt 
      ON m.id = mt.memos_id WHERE m.id = $2 AND m.users_id = $3;
    EOS
=end

    select_memo_query = <<~EOS
      SELECT DISTINCT
        m.id AS id,
        m.subject AS subject,
        m.content AS content
      FROM memos m JOIN memo_tag mt 
      ON m.id = mt.memos_id WHERE m.id = $1 AND m.users_id = $2;
    EOS

    select_memo_rslt = @conn.exec(select_memo_query, [memo_id, user_id])

    memos = []
    select_memo_rslt.each do |row|
      memos.push(row)
    end

    select_tags_query = <<~EOS
      SELECT t.id, t.name 
      FROM tags t
      JOIN memo_tag mt
      ON t.id = mt.tags_id
      WHERE mt.memos_id = $1
    EOS

    select_tags_rslt = @conn.exec(select_tags_query, [memo_id])

    tags = []
    select_tags_rslt.each do |row|
      tags.push(row)
    end

    return memos, tags
  end

  # upsert..できないか？シーケンスがネック
  # TODO: トランザクションはる
  # TODO: 関数に分割する
  def update(args)
    rslt = ''

    if args['memo_id'].empty?
      # メモ新規作成
      q = <<~EOS
        INSERT INTO memos(subject, content, users_id)
        VALUES($1, $2, $3)
        RETURNING id
      EOS

      rslt = @conn.exec(q, [
        args['subject'],
        args['content'],
        args['user_id']
      ])
    else
      # メモ編集
      q = <<~EOS
        UPDATE memos SET subject=$1, content=$2
        WHERE id=$3 AND users_id=$4
        RETURNING id
      EOS

      rslt = @conn.exec(q, [
        args['subject'],
        args['content'],
        args['memo_id'],
        args['user_id']
      ])
    end

    rslt[0]['id']
  end

  def fetch_all_tags_of_user(user_id)
    q = 'SELECT id, name FROM tags WHERE users_id = $1 ORDER BY id'

    rslt = @conn.exec(q, [user_id])

    rows = []
    rslt.each do |row|
      rows.push(row)
    end
    rows
  end
end
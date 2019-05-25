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
    q = 'SELECT id, subject FROM memos WHERE users_id=$1'
    rslt = @conn.exec(q, [user_id])

    rows = []
    rslt.each do |row|
      rows.push(row)
    end
    rows
  end

  def detail(memo_id, user_id)
    #q = 'SELECT id, subject, content FROM memos WHERE id=$1 AND users_id=$2'
    
    # メモ詳細画面にタグ情報を出力&メモ編集画面にタグ情報を渡すため.後日画面実装して確かめる
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
        ) AS tag_names 
      FROM memos m JOIN memo_tag mt 
      ON m.id = mt.memos_id WHERE m.id = $1 AND m.users_id = $2;
    EOS

    rslt = @conn.exec(q, [memo_id, user_id])

    rows = []
    rslt.each do |row|
      rows.push(row)
    end
    rows
  end

  # upsert..できないか？シーケンスがネック
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
end
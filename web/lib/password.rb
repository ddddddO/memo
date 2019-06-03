require 'digest/sha2'

module Password
  # パスワード＋Salt(ユーザー名)のハッシュ化*n(ストレッチング)結果を返却
  # ref: https://glorificatio.org/archives/3417
  def gen_secure_password(password, salt, n)
    return stretching(password, salt, n)
  end

  module_function :gen_secure_password

  def stretching(password, salt, n)
    secure_password = password + salt

    n.times do
      secure_password = Digest::SHA256.hexdigest(secure_password)
    end

    return secure_password
  end

  module_function :stretching
end

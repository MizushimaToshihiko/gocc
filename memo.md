
#### 【テストをCで書き直す】
 'funcs'と各テストケースを一つのCファイルにして、それをcompile関数に渡す
 それで出来たasmをgccに渡して実行ファイルにする
 以下を止める
 ・funcs_file作成
 ・1回1回テストケース毎にCファイルを作る => 一回だけになる


#### 【tokenize関数が\"を認識しない原因】
 <= string型の変数に入れた時点で\の部分が消えてしまうから？
 <= 'userInput'に入れるときに[]runeで入れる?
  ⇒　[]runeに変更

#### 【test実行時の出力がおかしい】
 <= 本家ではコンパイラがcなので、gcc -o $@ $(OBJS) -std=c11 -g -static
    でgccコンパイラから「printf」をstaticリンクで認識している
 => 本家と同じようにgccとstaticリンクする？どのように？
 　 <= gcc に渡すときに-cを付けなければリンクされるみたい
 　　　https://tanakamura.github.io/pllp/docs/linker.html なので問題なくリンクされている
 
  => ~~cコンパイラからgoコンパイラに切り替える？~~
  
  
  => 2021/11/5原因判明:readStringLiteral関数内の処理でtok.ContLenを1少なく設定したため、'\0'を読み込まず、'.byte 0'をasm.sに書き込んでいなかった初歩的ミスが原因だった
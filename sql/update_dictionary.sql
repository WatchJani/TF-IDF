INSERT INTO tag.dictionary (word, tokenization)
VALUES ('word', 1)
ON CONFLICT (word) DO UPDATE
SET tokenization = dictionary.tokenization + 1;
# Accounts.

We have three classess of users;
- anonymous, no log in required; read only;
- registered users; can read, can post annotations and comments;
- administrators; can do everything, uploading documents, adding metadata, add remove tags, delete stuff.

To become an admin, sign in at the site, (choose a good/strong password)
- run this on the server:

    mongo << EOF
    use nulpunt
    db.accounts.update({"username": "Bruce"}, {"$set": {"admin": true}})
    EOF

Change the name if you are not called Bruce.

Taking away admin rights is left as an excercise for the reader.
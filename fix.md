


  |-------------------------------|-----------------------------------------------|------------------------|
  | Unrestricted NOPASSWD sudo    | roles/bootstrap/tasks/wordsail-user.yml:32    | Full system compromise |
  | Plaintext DB credentials      | roles/database/templates/*.j2                 | Credential exposure    |
  | *.*:ALL,GRANT DB privileges   | roles/database/tasks/main.yml:57              | Database compromise    |
  | Excessive ignore_errors: true | roles/website/tasks/wordpress.yml:20,34,41,48 | Silent failures        |



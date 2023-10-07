
# Makromusic Case

# Başlangıç
Projeyi başlatmak ve gereksinimleri karşılamak için aşağıdaki adımları izleyin.

# Önkoşullar
Bu projeyi başlatmak için aşağıdaki önkoşulları karşılamalısınız:

- Docker kurulu olmalıdır.
- Kurulum
- Bu projeyi Git ile klonlayın:


git clone https://github.com/yigitsilto/faceDetection.git
Proje dizinine gidin:

 # SH Hakkında
runProtosAndProject.sh komutunu çalıştırarak Docker konteynerlerini başlatın ve proto çıktısını hazırlayın:

- sh runProtosAndProject.sh

# example_base64_image.json
Bu dosya, yüklenecek bir base64 kodlu bir JSON dosyasını içerir. Bu dosya, "upload image" fonksiyonunun örnek bir veri girişi olarak kullanılabilir.

# credentials.json
Bu dosya, Vision API'ye istek göndermek için gerekli kimlik bilgilerini içerir. Bu kimlik bilgileri, projenin otomatik olarak ayarlandığı ve kullanıldığı yerdir.

# Proje Adımları
# Upload Image
- Upload Image GRPC fonskiyonu yüklemek istediğiniz resimin base64 formatına göre diske yazar ve veritabanı (Postgres) e yazar.
- Kafka producer kodu çağrılır ve belirli bir topic üzerinden veriyi gönderir
- Consumer bu topiği dinler ve bir veri geldiğinde vision api servisini çağırır google apisini çağırır ve face detection işlemlerini yapıp veritabanına yazar. Aynı zamanda google dan dönen bütün json verisini de redise yazar

# GetImageDetail
- Get image detail fonksiyonu redisten google tarafından dönen bütün face detection json bilgisini belirli bir id değerine göre iletir.

# GetImageFeed
- Bu fonksiyon google tarafından sağlanan veri değerleriyle birlikte(Joy, Anger vs) resmin id, oluşma tarihi bilgilerininin olduğu bir listeyi sayfalama yaparak ve created at e göre desc sırasına sokarak döner

# UpdateImage
- Bu fonksiyon gönderilen id değerine göre kafkaya giderek face detection apisini tekrar çağırıp bilgileri günceller




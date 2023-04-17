class Asset{
  String UniName;
  String ID;
  String Name;
  String Major;
  int Credit;

  Asset(
      this.UniName,
      this.ID,
      this.Name,
      this.Major,
      this.Credit,
  );

  factory Asset.fromJson(Map<String,dynamic>json)=> Asset(
      json["uniName"],
      json["ID"],
      json["name"],
      json["major"],
      json["credit"],
  );

  Map<String,dynamic> toJson()=>{
    "uniName":UniName.toString(),
    "ID":ID.toString(),
    "name":Name.toString(),
    "major":Major.toString(),
    "credit":Credit.toString(),
  };
}
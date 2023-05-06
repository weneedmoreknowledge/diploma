class Asset{
  String University;
  String ID;
  String StudentName;
  int Credit;
  bool GradStatus;

  Asset(
      this.University,
      this.ID,
      this.StudentName,
      this.Credit,
      this.GradStatus,
  );

  factory Asset.fromJson(Map<String,dynamic>json)=> Asset(
    json['University'],
    json['ID'],
    json['StudentName'],
    json['Credit'],
    json['GradStatus'],
  );

  Map<String,dynamic> toJson()=>{
    'University':University.toString(),
    'ID':ID.toString(),
    'StudentName':StudentName.toString(),
    'Credit':Credit.toString(),
    'GradStatus':GradStatus.toString(),
  };
}
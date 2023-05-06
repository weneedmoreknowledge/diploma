import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'package:http/http.dart' as http;
import 'package:get/get.dart';
import 'package:fluttertoast/fluttertoast.dart';

import '../utils/API.dart';
import '../format/assets.dart';


class HomePage extends StatefulWidget {
  final String getUrlString;
  const HomePage({
    Key? key,
    required this.getUrlString
  }) : super(key: key);

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  
  Future getAllNow()async{
    List<Asset> assets=[];
    try{
      var res= await http.get(
        Uri.parse(widget.getUrlString),
      );
      if(res.statusCode==200){
        final resBodyOfAsset=jsonDecode(res.body);
        late int len = resBodyOfAsset.length;
        for(var i=0; i < len; i++){
          assets.add(Asset.fromJson(resBodyOfAsset[i]));
        }
        assets;
      }else{
        Fluttertoast.showToast(msg: 'Wrong information. Try again');
      }
    }catch(e){
      Fluttertoast.showToast(msg: e.toString());
      throw(e.toString());
    }
    return assets;
  }

  Future<void> refreshData() async {
    setState(() {
      getAllNow();
    });
  }

  @override
  void initState() {
    super.initState();
    getAllNow();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Hello World"),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: refreshData,
          ),
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: (){
              Get.to(()=>CreatePage());
            },
          ),
          IconButton(
            icon: const Icon(Icons.update),
            onPressed: (){
              Get.to(()=>UpdatePage());
            },
          ),
        ],
      ),
      body: FutureBuilder(
        future: getAllNow(),
        builder: (BuildContext context, AsyncSnapshot snapshot){
          Future.delayed(const Duration(seconds: 3));
          if(snapshot.hasError){
            print(snapshot.error.toString());
          }
          return Container(
            padding: const EdgeInsets.symmetric(vertical: 16),
            color: Colors.white,
            child: snapshot.hasData? ListView.builder(
              itemCount: snapshot.data.length,
                itemBuilder: (context,index){
                  return ListTile(
                    tileColor: Colors.blue,
                    leading: Text(snapshot.data[index].University),
                    title: Text(snapshot.data[index].StudentName),
                    subtitle: Text(snapshot.data[index].ID),
                    trailing: Text(snapshot.data[index].GradStatus.toString()),
                  );
                }
            ):const Center(child: CircularProgressIndicator(),),
          );
        }
      ),
    );
  }
}


class CreatePage extends StatefulWidget {
  const CreatePage({Key? key}) : super(key: key);

  @override
  State<CreatePage> createState() => _CreatePageState();
}

class _CreatePageState extends State<CreatePage> {
  final TextEditingController _controllerUniversity = TextEditingController();
  final TextEditingController _controllerID = TextEditingController();
  final TextEditingController _controllerStudentName = TextEditingController();

  @override
  void dispose() {
    _controllerUniversity.dispose();
    _controllerID.dispose();
    _controllerStudentName.dispose();
    super.dispose();
  }

  void submitAsset()async{
    final Map<String, String> requestData = {
      'university': _controllerUniversity.text,
      'id': _controllerID.text.trim(),
      'studentName': _controllerStudentName.text,
    };
    try{
      var res= await http.put(
        Uri.parse(API.createUrl),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(requestData),
      );
      if(res.statusCode==200){
        print('Data posted successfully');
        print(res.body);
        Get.back();
      }else{
        Fluttertoast.showToast(msg: '${res.statusCode}');
      }
    }catch(e){
      Fluttertoast.showToast(msg: e.toString());
      throw(e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          leading: IconButton(
            icon: const Icon(Icons.arrow_back_ios_sharp),
            onPressed: (){
              Get.back();
            },
          ),
          title: const Text("Create Page"),
        ),
        body: Column(
          children: [
            const SizedBox(height: 16,),
            TextFormField(
              controller: _controllerUniversity,
              decoration: const InputDecoration(
                  hintText: "Enter the University",
                  labelText: "University",
                  labelStyle: TextStyle(
                      color: Colors.black,
                  ),
                  hintStyle: TextStyle(
                      color: Colors.black
                  ),
                  fillColor:Colors.transparent,
                  filled: true,
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(
                        color: Colors.black
                    ),
                  )
              ),
            ),
            const SizedBox(height: 32,),
            TextFormField(
              controller: _controllerID,
              decoration: const InputDecoration(
                  hintText: 'Enter the ID',
                  labelText: 'ID',
                  labelStyle: TextStyle(
                      color: Colors.black
                  ),
                  hintStyle: TextStyle(
                      color: Colors.black
                  ),
                  fillColor:Colors.transparent,
                  filled: true,
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(
                        color: Colors.black
                    ),
                  )
              ),
            ),
            const SizedBox(height: 32,),
            TextFormField(
              controller: _controllerStudentName,
              decoration: const InputDecoration(
                  hintText: 'Enter the StudentName',
                  labelText: 'StudentName',
                  labelStyle: TextStyle(
                      color: Colors.black
                  ),
                  hintStyle: TextStyle(
                      color: Colors.black
                  ),
                  fillColor:Colors.transparent,
                  filled: true,
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(
                        color: Colors.black
                    ),
                  )
              ),
            ),
            ElevatedButton(
              onPressed: (){
                submitAsset();
              },
              child: const Text("Submit"),
            ),
          ],
        ),
      ),
    );
  }
}


class UpdatePage extends StatefulWidget {
  const UpdatePage({Key? key}) : super(key: key);

  @override
  State<UpdatePage> createState() => _UpdatePageState();
}

class _UpdatePageState extends State<UpdatePage> {
  final TextEditingController _controllerID = TextEditingController();
  final TextEditingController _controllerCredit = TextEditingController();
  bool isGraduated = false;

  @override
  void dispose() {
    _controllerID.dispose();
    _controllerCredit.dispose();
    super.dispose();
  }

  void submitAsset()async{
    final Map<String, String> requestData = {
      'id': _controllerID.text.trim(),
      'credit': _controllerCredit.text.trim(),
      'gradStatus': isGraduated.toString(),
    };
    try{
      var res= await http.put(
        Uri.parse(API.updateUrl),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(requestData),
      );
      if(res.statusCode==200){
        print('Data posted successfully');
        print(res.body);
        Get.back();
      }else{
        Fluttertoast.showToast(msg: '${res.statusCode}');
      }
    }catch(e){
      Fluttertoast.showToast(msg: e.toString());
      throw(e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          leading: IconButton(
            icon: const Icon(Icons.arrow_back_ios_sharp),
            onPressed: (){
              Get.back();
            },
          ),
          title: const Text("Update Page"),
        ),
        body: Column(
          children: [
            const SizedBox(height: 16,),
            TextFormField(
              controller: _controllerID,
              decoration: const InputDecoration(
                  hintText: "ID",
                  labelText: "Which ID do you want to update?",
                  labelStyle: TextStyle(
                    color: Colors.black,
                  ),
                  hintStyle: TextStyle(
                      color: Colors.black
                  ),
                  fillColor:Colors.transparent,
                  filled: true,
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(
                        color: Colors.black
                    ),
                  )
              ),
            ),
            const SizedBox(height: 32,),
            TextFormField(
              controller: _controllerCredit,
              keyboardType: TextInputType.number,
              inputFormatters: [
                FilteringTextInputFormatter.allow(RegExp(r'^[0-9]+$')),
              ],
              decoration: const InputDecoration(
                  hintText: 'Credit',
                  labelText: 'Enter the Credit',
                  labelStyle: TextStyle(
                      color: Colors.black
                  ),
                  hintStyle: TextStyle(
                      color: Colors.black
                  ),
                  fillColor:Colors.transparent,
                  filled: true,
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(
                        color: Colors.black
                    ),
                  )
              ),
            ),
            Switch(
                value: isGraduated,
                onChanged: (value) {
                  setState(() {
                    isGraduated = value;
                  });
                },
            ),
            ElevatedButton(
              onPressed: (){
                submitAsset();
              },
              child: const Text("Submit"),
            ),
          ],
        ),
      ),
    );
  }
}
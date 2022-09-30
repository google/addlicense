component persistent="true" table="todo"{
    property name="id" fieldtype="id" generator="native"; 
    property name="title" ormtype="string";
    property name="updated" ormtype="timestamp";
    property name="completed" ormtype="timestamp";

    public boolean function isComplete(){
        if (len(this.getCompleted()) GT 0 ){
            return true;
        } 
        return false;
    }


}
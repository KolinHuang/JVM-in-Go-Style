public class GaussTest{
	public static int staticVar;
	public int instanceVar;
	public static void main(String[] args){
		int x = 32768;
		GaussTest gt = new GaussTest();
		GaussTest.staticVar = x;
		x = GaussTest.staticVar;
		x = gt.instanceVar;
		gt.instanceVar = x;
		GaussTest gt2 = gt;
		if(gt2 instanceof GaussTest){
			gt = (GaussTest) gt2;
			System.out.println(gt.instanceVar);
		}
	}
}